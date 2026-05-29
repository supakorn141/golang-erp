package domain

import "time"

type User struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}

type UserUsecase interface {
	Register(name, email, password string) (*User, error)
	Login(email, password string) (string, error)
}
