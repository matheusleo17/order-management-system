package worker

import (
	"encoding/json"
	"log"

	"order-management-system/internal/events"

	"github.com/streadway/amqp"
)

const maxRetries = 3

func StartWorker(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		go handleMessage(msg)
	}
}

func handleMessage(msg amqp.Delivery) {
	retryCount := getRetryCount(msg)

	if retryCount >= maxRetries {
		log.Printf("Max retries (%d) reached, sending to DLQ", maxRetries)
		msg.Nack(false, false)
		return
	}

	var event events.OrderCreatedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal message (retry %d/%d): %v", retryCount+1, maxRetries, err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Processing order: %s (attempt %d/%d)", event.OrderID, retryCount+1, maxRetries)

	if err := processOrder(event); err != nil {
		log.Printf("Failed to process order %s (attempt %d/%d): %v", event.OrderID, retryCount+1, maxRetries, err)
		msg.Nack(false, true)
		return
	}

	log.Printf("Order processed successfully: %s", event.OrderID)
	msg.Ack(false)
}

func processOrder(event events.OrderCreatedEvent) error {
	// TODO: implement real business logic here
	// Example: notify inventory, send confirmation email, etc.
	log.Printf("Processing order %s with total %.2f", event.OrderID, event.Total)
	return nil
}

// getRetryCount reads the explicit retry header set by the publisher.
// Falls back to x-death count from RabbitMQ for backwards compatibility.
func getRetryCount(msg amqp.Delivery) int {
	// Check explicit retry header first
	if count, ok := msg.Headers["x-retry-count"].(int32); ok {
		return int(count)
	}

	// Fallback: count via RabbitMQ x-death header
	if deaths, ok := msg.Headers["x-death"].([]interface{}); ok && len(deaths) > 0 {
		if death, ok := deaths[0].(amqp.Table); ok {
			if count, ok := death["count"].(int64); ok {
				return int(count)
			}
		}
	}

	return 0
}
