package main

import (
	"log"

	"github.com/supakorn141/golang-erp/internal/config"
	"github.com/supakorn141/golang-erp/internal/handlers"
	"github.com/supakorn141/golang-erp/internal/middleware"
	"github.com/supakorn141/golang-erp/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	database.Migrate(db)

	r := gin.Default()
	r.Use(middleware.CORS())

	handlers.Register(r, db, cfg)

	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
