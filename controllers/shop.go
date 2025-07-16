package controllers

import (
	"net/http"
	"strconv"
	"tradesman-api/config"
	"tradesman-api/middleware"
	"tradesman-api/models"

	"github.com/gin-gonic/gin"
)

type ShopController struct{}

type CreateShopRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
}

// @Summary Tüm Esnafları Listele
// @Description Aktif olan tüm esnafları listeler
// @Tags Shops
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /shops [get]
func (sc *ShopController) GetShops(c *gin.Context) {
	var shops []models.Shop
	if err := config.DB.Preload("User").Where("is_active = ?", true).Find(&shops).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Esnaflar getirilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shops": shops,
	})
}

// @Summary Esnaf Detayı
// @Description Belirli bir esnafın detaylarını getirir
// @Tags Shops
// @Produce json
// @Param id path int true "Esnaf ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /shops/{id} [get]
func (sc *ShopController) GetShop(c *gin.Context) {
	id := c.Param("id")

	var shop models.Shop
	if err := config.DB.Preload("User").Preload("Products", "is_active = ?", true).First(&shop, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Esnaf bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shop": shop,
	})
}

// @Summary Esnaf Oluştur
// @Description Yeni esnaf kaydı oluşturur (sadece shop rolündeki kullanıcılar)
// @Tags Shops
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shop body CreateShopRequest true "Esnaf bilgileri"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /shops [post]
func (sc *ShopController) CreateShop(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	// Sadece shop rolündeki kullanıcılar esnaf oluşturabilir
	if userRole != models.RoleShop {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sadece esnaf rolündeki kullanıcılar dükkan oluşturabilir"})
		return
	}

	// Kullanıcının zaten bir dükkânı var mı kontrol et
	var existingShop models.Shop
	if err := config.DB.Where("user_id = ?", userID).First(&existingShop).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zaten bir dükkanınız var"})
		return
	}

	var req CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shop := models.Shop{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Phone:       req.Phone,
		IsActive:    true,
	}

	if err := config.DB.Create(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dükkan oluşturulamadı"})
		return
	}

	// Shop'u user ile birlikte getir
	config.DB.Preload("User").First(&shop, shop.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Dükkan başarıyla oluşturuldu",
		"shop":    shop,
	})
}

// @Summary Esnaf Güncelle
// @Description Esnaf bilgilerini günceller
// @Tags Shops
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Esnaf ID"
// @Param shop body CreateShopRequest true "Güncellenecek esnaf bilgileri"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /shops/{id} [put]
func (sc *ShopController) UpdateShop(c *gin.Context) {
	userID := middleware.GetUserID(c)
	shopID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz dükkan ID"})
		return
	}

	var shop models.Shop
	if err := config.DB.First(&shop, shopID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dükkan bulunamadı"})
		return
	}

	// Sadece dükkan sahibi güncelleyebilir
	if shop.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu dükkanı güncelleme yetkiniz yok"})
		return
	}

	var req CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shop.Name = req.Name
	shop.Description = req.Description
	shop.Address = req.Address
	shop.Phone = req.Phone

	if err := config.DB.Save(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dükkan güncellenemedi"})
		return
	}

	config.DB.Preload("User").First(&shop, shop.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Dükkan başarıyla güncellendi",
		"shop":    shop,
	})
}

// @Summary Esnafın Ürünlerini Listele
// @Description Belirli bir esnafın ürünlerini listeler
// @Tags Shops
// @Produce json
// @Param id path int true "Esnaf ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /shops/{id}/products [get]
func (sc *ShopController) GetShopProducts(c *gin.Context) {
	shopID := c.Param("id")

	var shop models.Shop
	if err := config.DB.First(&shop, shopID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Esnaf bulunamadı"})
		return
	}

	var products []models.Product
	if err := config.DB.Where("shop_id = ? AND is_active = ?", shopID, true).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürünler getirilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shop_id":   shop.ID,
		"shop_name": shop.Name,
		"products":  products,
	})
}
