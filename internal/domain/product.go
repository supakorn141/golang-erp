package domain

import "time"

type Product struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name"`
	SKU       string     `json:"sku"`
	Price     float64    `json:"price"`
	Stock     int        `json:"stock"`
	Category  string     `json:"category"`
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
