package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint        `gorm:"not null" json:"user_id"`
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status     string      `gorm:"default:pending" json:"status"` // pending, confirmed, shipped, completed, cancelled
	TotalPrice float64     `json:"total_price"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"not null" json:"unit_price"`
}
