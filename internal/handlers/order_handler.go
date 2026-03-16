package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"order-management-system/internal/domain"
	"order-management-system/internal/service"
)

func CreateOrder(c *gin.Context) {

	var order domain.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	result := service.CreateOrder(order)

	c.JSON(http.StatusCreated, result)
}
