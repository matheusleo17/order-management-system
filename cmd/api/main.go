package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"order-management-system/internal/database"
	"order-management-system/internal/handlers"
	"order-management-system/internal/messaging"
	"order-management-system/internal/repository"
	"order-management-system/internal/service"
)

func main() {

	err := database.ConnectMongo()

	if err != nil {
		log.Fatal("Mongo Connection Failed", err)
	}
	router := gin.Default()

	publisher, err := messaging.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	repo := repository.NewOrderRepository()
	service := service.NewOrderService(repo, publisher)
	handler := handlers.NewOrderHandler(service)

	router.POST("/orders", handler.CreateOrder)
	router.GET("/orders", handler.GetOrders)
	router.GET("/orders/:id", handler.GetOrderById)
	router.PUT("/orders/:id", handler.UpdateOrder)
	router.DELETE("/oders/:id", handler.DeleteOrder)

	router.Run(":8080")
}
