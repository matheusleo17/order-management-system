package worker

import (
	"encoding/json"
	"errors"
	"log"

	"order-management-system/internal/events"

	"github.com/streadway/amqp"
)

func StartWorker(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		go handleMessage(msg)
	}
}

func handleMessage(msg amqp.Delivery) {
	var event events.OrderCreatedEvent

	retryCount := getRetryCount(msg)
	maxRetries := 3

	if retryCount >= maxRetries {
		log.Println("Enviando para DLQ:", retryCount)

		msg.Nack(false, false)
		return
	}

	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		log.Println("erro", err)
		msg.Nack(false, true)
		return
	}
	log.Println("Processando pedido:", event.OrderID)

	err = processOrder(event)
	if err != nil {
		log.Println("erro no processamento:", err)

		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func processOrder(event events.OrderCreatedEvent) error {
	if event.Total > 1000 {
		return errors.New("erro simulado")
	}

	log.Println("Processado com sucesso:", event.OrderID)
	return nil
}

func getRetryCount(msg amqp.Delivery) int {
	if deaths, ok := msg.Headers["x-death"].([]interface{}); ok && len(deaths) > 0 {
		death := deaths[0].(amqp.Table)
		if count, ok := death["count"].(int64); ok {
			return int(count)
		}
	}
	return 0
}
