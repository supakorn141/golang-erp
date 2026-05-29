package repository

import (
	"database/sql"

	"github.com/supakorn141/golang-erp/internal/domain"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindAll() ([]domain.Order, error) {
	query := `
		SELECT o.id, o.created_at, o.updated_at, o.user_id, o.status, o.total_price,
		       u.id, u.name, u.email, u.role
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE o.deleted_at IS NULL
		ORDER BY o.id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		o.User = &domain.User{}
		if err := rows.Scan(
			&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.UserID, &o.Status, &o.TotalPrice,
			&o.User.ID, &o.User.Name, &o.User.Email, &o.User.Role,
		); err != nil {
			return nil, err
		}
		o.Items, err = r.findItems(o.ID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *orderRepository) FindByID(id uint) (*domain.Order, error) {
	query := `
		SELECT o.id, o.created_at, o.updated_at, o.user_id, o.status, o.total_price,
		       u.id, u.name, u.email, u.role
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE o.id = $1 AND o.deleted_at IS NULL`

	o := &domain.Order{User: &domain.User{}}
	err := r.db.QueryRow(query, id).Scan(
		&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.UserID, &o.Status, &o.TotalPrice,
		&o.User.ID, &o.User.Name, &o.User.Email, &o.User.Role,
	)
	if err != nil {
		return nil, err
	}
	o.Items, err = r.findItems(o.ID)
	return o, err
}

func (r *orderRepository) findItems(orderID uint) ([]domain.OrderItem, error) {
	query := `
		SELECT oi.id, oi.created_at, oi.order_id, oi.product_id, oi.quantity, oi.unit_price,
		       p.id, p.name, p.sku, p.price
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		item.Product = &domain.Product{}
		if err := rows.Scan(
			&item.ID, &item.CreatedAt, &item.OrderID, &item.ProductID,
			&item.Quantity, &item.UnitPrice,
			&item.Product.ID, &item.Product.Name, &item.Product.SKU, &item.Product.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *orderRepository) UpdateStatus(id uint, status string) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, status, id)
	return err
}
