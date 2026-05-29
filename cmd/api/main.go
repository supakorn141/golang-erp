package main

import (
	"log"

	"github.com/supakorn141/golang-erp/internal/config"
	deliveryHttp "github.com/supakorn141/golang-erp/internal/delivery/http"
	"github.com/supakorn141/golang-erp/internal/repository"
	"github.com/supakorn141/golang-erp/internal/usecase"
	"github.com/supakorn141/golang-erp/pkg/database"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	database.Migrate(db)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Usecases
	userUC := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	productUC := usecase.NewProductUsecase(productRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, productRepo, db)

	// Router
	r := deliveryHttp.NewRouter(userUC, productUC, orderUC, cfg.JWTSecret)

	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
