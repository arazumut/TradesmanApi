package controllers

import (
	"net/http"
	"strconv"
	"tradesman-api/config"
	"tradesman-api/middleware"
	"tradesman-api/models"

	"github.com/gin-gonic/gin"
)

type ProductController struct{}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	ImageURL    string  `json:"image_url"`
}

// @Summary Tüm Ürünleri Listele
// @Description Aktif olan tüm ürünleri listeler
// @Tags Products
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func (pc *ProductController) GetProducts(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Preload("Shop").Where("is_active = ?", true).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürünler getirilemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// @Summary Ürün Detayı
// @Description Belirli bir ürünün detaylarını getirir
// @Tags Products
// @Produce json
// @Param id path int true "Ürün ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [get]
func (pc *ProductController) GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := config.DB.Preload("Shop.User").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// @Summary Ürün Oluştur
// @Description Yeni ürün oluşturur (sadece esnaflar)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body CreateProductRequest true "Ürün bilgileri"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /products [post]
func (pc *ProductController) CreateProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	// Sadece esnaflar ürün ekleyebilir
	if userRole != models.RoleShop {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sadece esnaflar ürün ekleyebilir"})
		return
	}

	// Kullanıcının dükkanını bul
	var shop models.Shop
	if err := config.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Önce bir dükkan oluşturmalısınız"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		ShopID:      shop.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün oluşturulamadı"})
		return
	}

	// Product'ı shop ile birlikte getir
	config.DB.Preload("Shop").First(&product, product.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ürün başarıyla oluşturuldu",
		"product": product,
	})
}

// @Summary Ürün Güncelle
// @Description Ürün bilgilerini günceller (sadece ürün sahibi)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Ürün ID"
// @Param product body CreateProductRequest true "Güncellenecek ürün bilgileri"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ürün ID"})
		return
	}

	var product models.Product
	if err := config.DB.Preload("Shop").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	// Sadece ürün sahibi güncelleyebilir
	if product.Shop.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu ürünü güncelleme yetkiniz yok"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.ImageURL = req.ImageURL

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün güncellenemedi"})
		return
	}

	config.DB.Preload("Shop").First(&product, product.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ürün başarıyla güncellendi",
		"product": product,
	})
}

// @Summary Ürün Sil
// @Description Ürünü siler (sadece ürün sahibi)
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Ürün ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [delete]
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ürün ID"})
		return
	}

	var product models.Product
	if err := config.DB.Preload("Shop").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ürün bulunamadı"})
		return
	}

	// Sadece ürün sahibi silebilir
	if product.Shop.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu ürünü silme yetkiniz yok"})
		return
	}

	// Soft delete
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün silinemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ürün başarıyla silindi",
	})
}
