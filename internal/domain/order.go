package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint        `gorm:"not null" json:"user_id"`
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status     string      `gorm:"default:pending" json:"status"`
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

type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type OrderRepository interface {
	FindAll() ([]Order, error)
	FindByID(id uint) (*Order, error)
	Create(order *Order) error
	UpdateStatus(id uint, status string) error
}

type OrderUsecase interface {
	GetAll() ([]Order, error)
	GetByID(id uint) (*Order, error)
	Create(userID uint, items []CreateOrderItemRequest) (*Order, error)
	UpdateStatus(id uint, status string) error
}
