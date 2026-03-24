package main

import (
	"log"

	"github.com/gin-gonic/gin"

	handlers "order-management-system/internal/delivery/http"
	"order-management-system/internal/infrastructure/database"
	"order-management-system/internal/infrastructure/database/messaging"
	"order-management-system/internal/infrastructure/repository"
	"order-management-system/internal/usecase"
)

func main() {

	err := database.ConnectMongo()
	if err != nil {
		log.Fatal("Mongo Connection Failed", err)
	}

	router := gin.Default()

	repo := repository.NewOrderRepositoryMongo()
	publisher, err := messaging.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/")

	createUC := usecase.NewCreateOrderUseCase(repo, publisher)
	getAllUC := usecase.NewGetOrdersUseCase(repo)
	getByIdUC := usecase.NewGetOrderByIdUseCase(repo)
	updateUC := usecase.NewUpdateOrderUseCase(repo)
	deleteUC := usecase.NewDeleteOrderUseCase(repo)

	handler := handlers.NewOrderHandler(
		createUC,
		getAllUC,
		getByIdUC,
		updateUC,
		deleteUC,
	)

	router.POST("/orders", handler.CreateOrder)
	router.GET("/orders", handler.GetOrders)
	router.GET("/orders/:id", handler.GetOrderById)
	router.PUT("/orders/:id", handler.UpdateOrder)
	router.DELETE("/orders/:id", handler.DeleteOrder)

	router.Run(":8080")
}
