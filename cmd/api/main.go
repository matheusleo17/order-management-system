package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	handlers "order-management-system/internal/delivery/http"
	"order-management-system/internal/infrastructure/database"
	"order-management-system/internal/infrastructure/messaging"
	"order-management-system/internal/infrastructure/repository"
	"order-management-system/internal/logger"
	"order-management-system/internal/usecase"
)

func main() {
	log := logger.New()
	defer log.Sync()
	if err := database.ConnectMongo(log); err != nil {
		log.Fatal("Mongo connection failed", zap.Error(err))
	}

	publisher, err := messaging.NewRabbitMQPublisher(log)
	if err != nil {
		log.Fatal("RabbitMQ connection failed", zap.Error(err))
	}

	repo := repository.NewOrderRepositoryMongo()

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

	router := gin.Default()

	router.POST("/orders", handler.CreateOrder)
	router.GET("/orders", handler.GetOrders)
	router.GET("/orders/:id", handler.GetOrderById)
	router.PUT("/orders/:id", handler.UpdateOrder)
	router.DELETE("/orders/:id", handler.DeleteOrder)
	log.Info("API starting", zap.String("port", "8080"))

	router.Run(":8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
