package repository

import (
	"database/sql"

	"github.com/supakorn141/golang-erp/internal/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]domain.Product, error) {
	query := `
		SELECT id, created_at, updated_at, name, sku, price, stock, category
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt,
			&p.Name, &p.SKU, &p.Price, &p.Stock, &p.Category); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *productRepository) FindByID(id uint) (*domain.Product, error) {
	query := `
		SELECT id, created_at, updated_at, name, sku, price, stock, category
		FROM products
		WHERE id = $1 AND deleted_at IS NULL`

	p := &domain.Product{}
	err := r.db.QueryRow(query, id).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt,
			&p.Name, &p.SKU, &p.Price, &p.Stock, &p.Category)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *productRepository) Create(p *domain.Product) error {
	query := `
		INSERT INTO products (name, sku, price, stock, category)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, p.Name, p.SKU, p.Price, p.Stock, p.Category).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *productRepository) Update(p *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, sku = $2, price = $3, stock = $4, category = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING updated_at`

	return r.db.QueryRow(query, p.Name, p.SKU, p.Price, p.Stock, p.Category, p.ID).
		Scan(&p.UpdatedAt)
}

func (r *productRepository) Delete(id uint) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, id)
	return err
}
