package handlers

import (
	"net/http"

	"github.com/supakorn141/golang-erp/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func (h *OrderHandler) List(c *gin.Context) {
	var orders []models.Order
	h.db.Preload("Items.Product").Preload("User").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) Get(c *gin.Context) {
	var order models.Order
	if err := h.db.Preload("Items.Product").Preload("User").First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

type createOrderRequest struct {
	Items []struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	order := models.Order{UserID: userID}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		var total float64
		for _, item := range req.Items {
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return err
			}
			if product.Stock < item.Quantity {
				return gorm.ErrInvalidData
			}
			tx.Model(&product).Update("stock", product.Stock-item.Quantity)
			orderItem := models.OrderItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				UnitPrice: product.Price,
			}
			order.Items = append(order.Items, orderItem)
			total += product.Price * float64(item.Quantity)
		}
		order.TotalPrice = total
		return tx.Create(&order).Error
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Model(&models.Order{}).Where("id = ?", c.Param("id")).Update("status", body.Status).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": body.Status})
}
