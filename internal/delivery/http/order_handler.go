package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"order-management-system/internal/domain"
)

// Use case interfaces — desacoplamento do handler das implementações concretas
type createOrderUseCase interface {
	Execute(ctx context.Context, order *domain.Order) error
}

type getOrdersUseCase interface {
	Execute(ctx context.Context) ([]domain.Order, error)
}

type getOrderByIdUseCase interface {
	Execute(ctx context.Context, id string) (domain.Order, error)
}

type updateOrderUseCase interface {
	Execute(ctx context.Context, id string, order *domain.Order) error
}

type deleteOrderUseCase interface {
	Execute(ctx context.Context, id string) error
}

type OrderHandler struct {
	createUC  createOrderUseCase
	getAllUC   getOrdersUseCase
	getByIdUC getOrderByIdUseCase
	updateUC  updateOrderUseCase
	deleteUC  deleteOrderUseCase
}

func NewOrderHandler(
	createUC createOrderUseCase,
	getAllUC getOrdersUseCase,
	getByIdUC getOrderByIdUseCase,
	updateUC updateOrderUseCase,
	deleteUC deleteOrderUseCase,
) *OrderHandler {
	return &OrderHandler{
		createUC:  createUC,
		getAllUC:   getAllUC,
		getByIdUC: getByIdUC,
		updateUC:  updateUC,
		deleteUC:  deleteUC,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order domain.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.createUC.Execute(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": order})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.getAllUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	id := c.Param("id")

	order, err := h.getByIdUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.updateUC.Execute(c.Request.Context(), id, &order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated successfully"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.deleteUC.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}
