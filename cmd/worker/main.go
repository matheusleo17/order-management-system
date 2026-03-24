package main

import (
	"log"

	"order-management-system/internal/infrastructure/messaging"
	"order-management-system/internal/worker"
)

func main() {
	consumer, err := messaging.NewRabbitMQConsumer()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}

	msgs, err := consumer.Consume("order.created")
	if err != nil {
		log.Fatal("Failed to start consuming: ", err)
	}

	log.Println("Worker running...")
	worker.StartWorker(msgs)
}
