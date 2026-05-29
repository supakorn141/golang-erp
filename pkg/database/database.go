package database

import (
	"fmt"
	"log"

	"github.com/supakorn141/golang-erp/internal/config"
	"github.com/supakorn141/golang-erp/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
