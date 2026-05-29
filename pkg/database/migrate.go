package database

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ,
			name       VARCHAR(255) NOT NULL,
			email      VARCHAR(255) NOT NULL UNIQUE,
			password   VARCHAR(255) NOT NULL,
			role       VARCHAR(50)  NOT NULL DEFAULT 'staff'
		)`,

		`CREATE TABLE IF NOT EXISTS products (
			id         SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ,
			name       VARCHAR(255) NOT NULL,
			sku        VARCHAR(100) NOT NULL UNIQUE,
			price      NUMERIC(15,2) NOT NULL DEFAULT 0,
			stock      INTEGER      NOT NULL DEFAULT 0,
			category   VARCHAR(100)
		)`,

		`CREATE TABLE IF NOT EXISTS orders (
			id          SERIAL PRIMARY KEY,
			created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			deleted_at  TIMESTAMPTZ,
			user_id     INTEGER      NOT NULL REFERENCES users(id),
			status      VARCHAR(50)  NOT NULL DEFAULT 'pending',
			total_price NUMERIC(15,2) NOT NULL DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS order_items (
			id         SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			order_id   INTEGER      NOT NULL REFERENCES orders(id),
			product_id INTEGER      NOT NULL REFERENCES products(id),
			quantity   INTEGER      NOT NULL,
			unit_price NUMERIC(15,2) NOT NULL
		)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			log.Fatalf("migration failed: %v\nSQL: %s", err, stmt)
		}
	}
	log.Println("migration complete")
}
