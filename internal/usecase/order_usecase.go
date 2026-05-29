package usecase

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/supakorn141/golang-erp/internal/domain"
)

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	db          *sql.DB
}

func NewOrderUsecase(orderRepo domain.OrderRepository, db *sql.DB) domain.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo, db: db}
}

func (u *orderUsecase) GetAll() ([]domain.Order, error) {
	return u.orderRepo.FindAll()
}

func (u *orderUsecase) GetByID(id uint) (*domain.Order, error) {
	return u.orderRepo.FindByID(id)
}

func (u *orderUsecase) Create(userID uint, items []domain.CreateOrderItemRequest) (*domain.Order, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	order := &domain.Order{UserID: userID, Status: "pending"}
	var total float64

	for _, item := range items {
		// ดึง product พร้อม lock row เพื่อป้องกัน race condition
		var product domain.Product
		row := tx.QueryRow(`
			SELECT id, name, price, stock
			FROM products
			WHERE id = $1 AND deleted_at IS NULL
			FOR UPDATE`, item.ProductID)

		if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s", product.Name)
		}

		// ลด stock
		if _, err := tx.Exec(`
			UPDATE products SET stock = stock - $1, updated_at = NOW()
			WHERE id = $2`, item.Quantity, product.ID); err != nil {
			return nil, err
		}

		order.Items = append(order.Items, domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
		})
		total += product.Price * float64(item.Quantity)
	}
	order.TotalPrice = total

	// สร้าง order
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, status, total_price)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`,
		order.UserID, order.Status, order.TotalPrice,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// สร้าง order items
	for i := range order.Items {
		err = tx.QueryRow(`
			INSERT INTO order_items (order_id, product_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at`,
			order.ID, order.Items[i].ProductID,
			order.Items[i].Quantity, order.Items[i].UnitPrice,
		).Scan(&order.Items[i].ID, &order.Items[i].CreatedAt)
		if err != nil {
			return nil, err
		}
		order.Items[i].OrderID = order.ID
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}
	return order, nil
}

func (u *orderUsecase) UpdateStatus(id uint, status string) error {
	return u.orderRepo.UpdateStatus(id, status)
}
