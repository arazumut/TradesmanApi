// @title Esnaf Yönetim Sistemi API
// @version 1.0
// @description Mahalledeki küçük esnafın ürünlerini yönetebileceği ve mobil uygulama üzerinden gelen müşteri siparişlerini takip edebileceği backend API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log"
	"tradesman-api/config"
	_ "tradesman-api/docs" // Swagger docs
	"tradesman-api/routes"
)

func main() {
	// Veritabanı bağlantısı
	config.InitDatabase()

	// Routes kurulumu
	r := routes.SetupRoutes()

	// Server başlat
	log.Println("🚀 Server başlatılıyor... http://localhost:8080")
	log.Println("📚 Swagger dokümantasyonu: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
}
