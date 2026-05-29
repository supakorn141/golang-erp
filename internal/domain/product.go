package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name     string  `gorm:"not null" json:"name"`
	SKU      string  `gorm:"uniqueIndex;not null" json:"sku"`
	Price    float64 `gorm:"not null" json:"price"`
	Stock    int     `gorm:"default:0" json:"stock"`
	Category string  `json:"category"`
}

type ProductRepository interface {
	FindAll() ([]Product, error)
	FindByID(id uint) (*Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id uint) error
}

type ProductUsecase interface {
	GetAll() ([]Product, error)
	GetByID(id uint) (*Product, error)
	Create(product *Product) error
	Update(id uint, input *Product) error
	Delete(id uint) error
}
