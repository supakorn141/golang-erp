package http

import (
	"github.com/gin-gonic/gin"
	"github.com/supakorn141/golang-erp/internal/delivery/http/middleware"
	"github.com/supakorn141/golang-erp/internal/domain"
)

func NewRouter(
	userUC domain.UserUsecase,
	productUC domain.ProductUsecase,
	orderUC domain.OrderUsecase,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")

	auth := NewAuthHandler(userUC)
	api.POST("/auth/register", auth.Register)
	api.POST("/auth/login", auth.Login)

	protected := api.Group("/")
	protected.Use(middleware.Auth(jwtSecret))

	product := NewProductHandler(productUC)
	protected.GET("/products", product.List)
	protected.GET("/products/:id", product.Get)
	protected.POST("/products", product.Create)
	protected.PUT("/products/:id", product.Update)
	protected.DELETE("/products/:id", product.Delete)

	order := NewOrderHandler(orderUC)
	protected.GET("/orders", order.List)
	protected.GET("/orders/:id", order.Get)
	protected.POST("/orders", order.Create)
	protected.PATCH("/orders/:id/status", order.UpdateStatus)

	return r
}
