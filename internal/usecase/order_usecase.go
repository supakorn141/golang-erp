package usecase

import (
	"errors"

	"github.com/supakorn141/golang-erp/internal/domain"
	"gorm.io/gorm"
)

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
	db          *gorm.DB
}

func NewOrderUsecase(orderRepo domain.OrderRepository, productRepo domain.ProductRepository, db *gorm.DB) domain.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo, productRepo: productRepo, db: db}
}

func (u *orderUsecase) GetAll() ([]domain.Order, error) {
	return u.orderRepo.FindAll()
}

func (u *orderUsecase) GetByID(id uint) (*domain.Order, error) {
	return u.orderRepo.FindByID(id)
}

func (u *orderUsecase) Create(userID uint, items []domain.CreateOrderItemRequest) (*domain.Order, error) {
	order := &domain.Order{UserID: userID}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		var total float64
		for _, item := range items {
			var product domain.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return err
			}
			if product.Stock < item.Quantity {
				return errors.New("insufficient stock for product: " + product.Name)
			}
			tx.Model(&product).Update("stock", product.Stock-item.Quantity)
			order.Items = append(order.Items, domain.OrderItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				UnitPrice: product.Price,
			})
			total += product.Price * float64(item.Quantity)
		}
		order.TotalPrice = total
		return tx.Create(order).Error
	})

	if err != nil {
		return nil, err
	}
	return order, nil
}

func (u *orderUsecase) UpdateStatus(id uint, status string) error {
	return u.orderRepo.UpdateStatus(id, status)
}
