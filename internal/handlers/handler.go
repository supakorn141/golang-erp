package handlers

import (
	"github.com/supakorn141/golang-erp/internal/config"
	"github.com/supakorn141/golang-erp/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	api := r.Group("/api/v1")

	// Public routes
	auth := &AuthHandler{db: db, secret: cfg.JWTSecret}
	api.POST("/auth/register", auth.Register)
	api.POST("/auth/login", auth.Login)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.Auth(cfg.JWTSecret))

	products := &ProductHandler{db: db}
	protected.GET("/products", products.List)
	protected.GET("/products/:id", products.Get)
	protected.POST("/products", products.Create)
	protected.PUT("/products/:id", products.Update)
	protected.DELETE("/products/:id", products.Delete)

	orders := &OrderHandler{db: db}
	protected.GET("/orders", orders.List)
	protected.GET("/orders/:id", orders.Get)
	protected.POST("/orders", orders.Create)
	protected.PATCH("/orders/:id/status", orders.UpdateStatus)
}
