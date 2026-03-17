package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"order-management-system/internal/domain"
	"order-management-system/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}
func (h *OrderHandler) CreateOrder(c *gin.Context) {

	var order domain.Order
	ctx := c.Request.Context()
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	result, err := h.service.CreateOrder(ctx, order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create order",
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}
func (h *OrderHandler) GetOrders(c *gin.Context) {
	ctx := c.Request.Context()

	orders, err := h.service.GetOrders(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get orders",
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}
func (h *OrderHandler) GetOrderById(c *gin.Context) {
	ctx := c.Param("id")

	orders, err := h.service.GetOrderById(c, ctx)

	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "failed to get orders",
		})
		return
	}
	c.JSON(http.StatusOK, orders)

}
