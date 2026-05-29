package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/supakorn141/golang-erp/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}
