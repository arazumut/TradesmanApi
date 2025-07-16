package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleShop     UserRole = "shop"     // Esnaf
	RoleCustomer UserRole = "customer" // Müşteri
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Phone     string         `json:"phone"`
	Role      UserRole       `json:"role" gorm:"type:varchar(20);default:'customer'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// İlişkiler
	Shop   *Shop   `json:"shop,omitempty" gorm:"foreignKey:UserID"`
	Orders []Order `json:"orders,omitempty" gorm:"foreignKey:UserID"`
}
