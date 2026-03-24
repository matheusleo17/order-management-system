package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"order-management-system/internal/domain"
	"order-management-system/internal/usecase"
)

type OrderHandler struct {
	createUC  *usecase.CreateOrderUseCase
	getAllUC  *usecase.GetOrdersUseCase
	getByIdUC *usecase.GetOrderByIdUseCase
	updateUC  *usecase.UpdateOrderUseCase
	deleteUC  *usecase.DeleteOrderUseCase
}

func NewOrderHandler(
	createUC *usecase.CreateOrderUseCase,
	getAllUC *usecase.GetOrdersUseCase,
	getByIdUC *usecase.GetOrderByIdUseCase,
	updateUC *usecase.UpdateOrderUseCase,
	deleteUC *usecase.DeleteOrderUseCase,
) *OrderHandler {
	return &OrderHandler{
		createUC:  createUC,
		getAllUC:  getAllUC,
		getByIdUC: getByIdUC,
		updateUC:  updateUC,
		deleteUC:  deleteUC,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order domain.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := h.createUC.Execute(c.Request.Context(), &order)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.getAllUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	id := c.Param("id")

	order, err := h.getByIdUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var order domain.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := h.updateUC.Execute(c.Request.Context(), id, &order)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "updated")
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	err := h.deleteUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, "deleted")
}
