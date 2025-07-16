package config

import (
	"log"
	"tradesman-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("tradesman.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}

	// Auto Migration
	err = DB.AutoMigrate(
		&models.User{},
		&models.Shop{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Veritabanı migrasyonu başarısız:", err)
	}

	log.Println("✅ Veritabanı başarıyla bağlandı ve migrate edildi!")
}
