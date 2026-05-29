package repository

import (
	"database/sql"

	"github.com/supakorn141/golang-erp/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, user.Name, user.Email, user.Password, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, created_at, updated_at, name, email, password, role
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt,
			&user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	query := `
		SELECT id, created_at, updated_at, name, email, password, role
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt,
			&user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
