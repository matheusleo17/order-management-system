package main

import (
	"log"

	"order-management-system/internal/messaging"
	"order-management-system/internal/worker"
)

func main() {
	consumer, err := messaging.NewRabbitMQConsumer("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := consumer.Consume("order.created")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Worker rodando...")

	worker.StartWorker(msgs)
}
