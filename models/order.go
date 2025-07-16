package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // Beklemede
	OrderStatusConfirmed OrderStatus = "confirmed" // Onaylandı
	OrderStatusPreparing OrderStatus = "preparing" // Hazırlanıyor
	OrderStatusReady     OrderStatus = "ready"     // Hazır
	OrderStatusDelivered OrderStatus = "delivered" // Teslim edildi
	OrderStatusCancelled OrderStatus = "cancelled" // İptal edildi
)

type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	ShopID      uint           `json:"shop_id" gorm:"not null;index"`
	TotalAmount float64        `json:"total_amount" gorm:"not null"`
	Status      OrderStatus    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Note        string         `json:"note"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// İlişkiler
	User       User        `json:"user" gorm:"foreignKey:UserID"`
	Shop       Shop        `json:"shop" gorm:"foreignKey:ShopID"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null;index"`
	ProductID uint      `json:"product_id" gorm:"not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"` // Sipariş anındaki fiyat
	CreatedAt time.Time `json:"created_at"`

	// İlişkiler
	Order   Order   `json:"order" gorm:"foreignKey:OrderID"`
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}
