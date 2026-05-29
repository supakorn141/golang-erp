package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name     string  `gorm:"not null" json:"name"`
	SKU      string  `gorm:"uniqueIndex;not null" json:"sku"`
	Price    float64 `gorm:"not null" json:"price"`
	Stock    int     `gorm:"default:0" json:"stock"`
	Category string  `json:"category"`
}
