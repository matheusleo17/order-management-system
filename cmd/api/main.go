package main

import (
	"github.com/gin-gonic/gin"

	"order-management-system/internal/handlers"
)

func main() {

	router := gin.Default()

	router.POST("/orders", handlers.CreateOrder)

	router.Run(":8080")
}
