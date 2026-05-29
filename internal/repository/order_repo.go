package repository

import (
	"github.com/supakorn141/golang-erp/internal/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindAll() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Items.Product").Preload("User").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items.Product").Preload("User").First(&order, id).Error
	return &order, err
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}
