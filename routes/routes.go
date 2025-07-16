package routes

import (
	"tradesman-api/controllers"
	"tradesman-api/middleware"
	"tradesman-api/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "🧾 Esnaf Yönetim Sistemi API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Controllers
	authController := &controllers.AuthController{}
	shopController := &controllers.ShopController{}
	productController := &controllers.ProductController{}
	orderController := &controllers.OrderController{}

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Public shop and product routes (for customers to browse)
	public := r.Group("/")
	{
		public.GET("/shops", shopController.GetShops)
		public.GET("/shops/:id", shopController.GetShop)
		public.GET("/shops/:id/products", shopController.GetShopProducts)
		public.GET("/products", productController.GetProducts)
		public.GET("/products/:id", productController.GetProduct)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth routes
		protected.GET("/auth/me", authController.Me)

		// Shop management (only for shop role)
		shopRoutes := protected.Group("/shops")
		shopRoutes.Use(middleware.RequireRole(models.RoleShop))
		{
			shopRoutes.POST("", shopController.CreateShop)
			shopRoutes.PUT("/:id", shopController.UpdateShop)
		}

		// Product management (only for shop role)
		productRoutes := protected.Group("/products")
		{
			productRoutes.POST("", middleware.RequireRole(models.RoleShop), productController.CreateProduct)
			productRoutes.PUT("/:id", middleware.RequireRole(models.RoleShop), productController.UpdateProduct)
			productRoutes.DELETE("/:id", middleware.RequireRole(models.RoleShop), productController.DeleteProduct)
		}

		// Order management
		orderRoutes := protected.Group("/orders")
		{
			orderRoutes.POST("", middleware.RequireRole(models.RoleCustomer), orderController.CreateOrder)
			orderRoutes.GET("", orderController.GetMyOrders)
			orderRoutes.GET("/:id", orderController.GetOrder)
			orderRoutes.PUT("/:id/status", middleware.RequireRole(models.RoleShop), orderController.UpdateOrderStatus)
		}
	}

	return r
}
