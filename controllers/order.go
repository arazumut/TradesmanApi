package controllers

import (
	"net/http"
	"strconv"
	"tradesman-api/config"
	"tradesman-api/middleware"
	"tradesman-api/models"

	"github.com/gin-gonic/gin"
)

type OrderController struct{}

type CreateOrderRequest struct {
	ShopID uint        `json:"shop_id" binding:"required"`
	Items  []OrderItem `json:"items" binding:"required,min=1"`
	Note   string      `json:"note"`
}

type OrderItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// @Summary Sipariş Oluştur
// @Description Yeni sipariş oluşturur (sadece müşteriler)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body CreateOrderRequest true "Sipariş bilgileri"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /orders [post]
func (oc *OrderController) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	// Sadece müşteriler sipariş verebilir
	if userRole != models.RoleCustomer {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sadece müşteriler sipariş verebilir"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Shop kontrolü
	var shop models.Shop
	if err := config.DB.First(&shop, req.ShopID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dükkan bulunamadı"})
		return
	}

	if !shop.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dükkan aktif değil"})
		return
	}

	// Transaction başlat
	tx := config.DB.Begin()

	var totalAmount float64 = 0
	var orderItems []models.OrderItem

	// Her ürün için kontrol yap
	for _, item := range req.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ürün bulunamadı: " + strconv.Itoa(int(item.ProductID))})
			return
		}

		if product.ShopID != req.ShopID {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ürün bu dükkanın değil"})
			return
		}

		if !product.IsActive {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ürün aktif değil: " + product.Name})
			return
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Yetersiz stok. Mevcut: " + strconv.Itoa(product.Stock) + ", İstenen: " + strconv.Itoa(item.Quantity),
			})
			return
		}

		// Stok güncelle
		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Stok güncellenemedi"})
			return
		}

		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Sipariş anındaki fiyat
		}
		orderItems = append(orderItems, orderItem)
		totalAmount += product.Price * float64(item.Quantity)
	}

	// Order oluştur
	order := models.Order{
		UserID:      userID,
		ShopID:      req.ShopID,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusPending,
		Note:        req.Note,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sipariş oluşturulamadı"})
		return
	}

	// Order items oluştur
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := tx.Create(&orderItems[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Sipariş kalemleri oluşturulamadı"})
			return
		}
	}

	// Transaction commit
	tx.Commit()

	// Order'ı ilişkilerle birlikte getir
	config.DB.Preload("Shop").Preload("OrderItems.Product").First(&order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sipariş başarıyla oluşturuldu",
		"order":   order,
	})
}

// @Summary Kullanıcının Siparişlerini Listele
// @Description Mevcut kullanıcının siparişlerini listeler
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /orders [get]
func (oc *OrderController) GetMyOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	var orders []models.Order

	if userRole == models.RoleCustomer {
		// Müşteriler sadece kendi siparişlerini görebilir
		config.DB.Preload("Shop").Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders)
	} else if userRole == models.RoleShop {
		// Esnaflar sadece kendi dükkanlarına gelen siparişleri görebilir
		var shop models.Shop
		if err := config.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dükkan bulunamadı"})
			return
		}
		config.DB.Preload("User").Preload("OrderItems.Product").Where("shop_id = ?", shop.ID).Find(&orders)
	} else {
		// Admin tüm siparişleri görebilir
		config.DB.Preload("User").Preload("Shop").Preload("OrderItems.Product").Find(&orders)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// @Summary Sipariş Detayı
// @Description Belirli bir siparişin detaylarını getirir
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sipariş ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (oc *OrderController) GetOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.Preload("User").Preload("Shop").Preload("OrderItems.Product").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sipariş bulunamadı"})
		return
	}

	// Yetki kontrolü
	canAccess := false
	if userRole == models.RoleAdmin {
		canAccess = true
	} else if userRole == models.RoleCustomer && order.UserID == userID {
		canAccess = true
	} else if userRole == models.RoleShop {
		var shop models.Shop
		if err := config.DB.Where("user_id = ?", userID).First(&shop).Error; err == nil {
			if order.ShopID == shop.ID {
				canAccess = true
			}
		}
	}

	if !canAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu siparişi görme yetkiniz yok"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

// @Summary Sipariş Durumu Güncelle
// @Description Sipariş durumunu günceller (sadece esnaflar)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sipariş ID"
// @Param status body map[string]string true "Yeni durum"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/status [put]
func (oc *OrderController) UpdateOrderStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)
	orderID := c.Param("id")

	// Sadece esnaflar sipariş durumu güncelleyebilir
	if userRole != models.RoleShop {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sadece esnaflar sipariş durumu güncelleyebilir"})
		return
	}

	var order models.Order
	if err := config.DB.Preload("Shop").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sipariş bulunamadı"})
		return
	}

	// Sipariş bu esnafın mı kontrol et
	if order.Shop.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu siparişi güncelleme yetkiniz yok"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Geçerli durum kontrolü
	validStatuses := []models.OrderStatus{
		models.OrderStatusPending,
		models.OrderStatusConfirmed,
		models.OrderStatusPreparing,
		models.OrderStatusReady,
		models.OrderStatusDelivered,
		models.OrderStatusCancelled,
	}

	isValidStatus := false
	for _, status := range validStatuses {
		if string(status) == req.Status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz sipariş durumu"})
		return
	}

	order.Status = models.OrderStatus(req.Status)
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sipariş durumu güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sipariş durumu güncellendi",
		"order":   order,
	})
}
