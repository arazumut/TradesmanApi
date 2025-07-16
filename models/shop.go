package models

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;uniqueIndex"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Address     string         `json:"address"`
	Phone       string         `json:"phone"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// İlişkiler
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Products []Product `json:"products,omitempty" gorm:"foreignKey:ShopID"`
	Orders   []Order   `json:"orders,omitempty" gorm:"foreignKey:ShopID"`
}
