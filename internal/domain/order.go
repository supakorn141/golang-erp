package domain

import "time"

type Order struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-"`
	UserID     uint       `json:"user_id"`
	User       *User      `json:"user,omitempty"`
	Status     string     `json:"status"`
	TotalPrice float64    `json:"total_price"`
	Items      []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Product   *Product  `json:"product,omitempty"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity"   binding:"required,min=1"`
}

type OrderRepository interface {
	FindAll() ([]Order, error)
	FindByID(id uint) (*Order, error)
	UpdateStatus(id uint, status string) error
}

type OrderUsecase interface {
	GetAll() ([]Order, error)
	GetByID(id uint) (*Order, error)
	Create(userID uint, items []CreateOrderItemRequest) (*Order, error)
	UpdateStatus(id uint, status string) error
}
