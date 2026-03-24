package main

import (
	"order-management-system/internal/infrastructure/messaging"
	"order-management-system/internal/logger"
	"order-management-system/internal/worker"

	"go.uber.org/zap"
)

func main() {
	log := logger.New()
	defer log.Sync()
	consumer, err := messaging.NewRabbitMQConsumer()

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", zap.Error(err))
	}

	msgs, err := consumer.Consume("order.created")
	if err != nil {
		log.Fatal("Failed to start consuming: ", zap.Error(err))
	}

	log.Info("Worker running...")
	worker.StartWorker(msgs, log)
}
