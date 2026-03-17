package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"order-management-system/internal/database"
	"order-management-system/internal/handlers"
	"order-management-system/internal/repository"
	"order-management-system/internal/service"
)

func main() {

	err := database.ConnectMongo()

	if err != nil {
		log.Fatal("Mongo Connection Failed", err)
	}
	router := gin.Default()

	repo := repository.NewOrderRepository()
	service := service.NewOrderService(repo)
	handler := handlers.NewOrderHandler(service)

	router.POST("/orders", handler.CreateOrder)
	router.GET("/orders", handler.GetOrders)
	router.GET("/orders/:id", handler.GetOrderById)

	router.Run(":8080")
}
