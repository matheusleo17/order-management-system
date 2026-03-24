package worker

import (
	"encoding/json"

	"order-management-system/internal/events"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const maxRetries = 3

func StartWorker(msgs <-chan amqp.Delivery, log *zap.Logger) {
	for msg := range msgs {
		go handleMessage(msg, log)
	}
}

func handleMessage(msg amqp.Delivery, log *zap.Logger) {
	retryCount := getRetryCount(msg)

	if retryCount >= maxRetries {
		log.Warn("Max retries reached, sending to DLQ",
			zap.Int("maxRetries", maxRetries))
		msg.Nack(false, false)
		return
	}

	var event events.OrderCreatedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Warn("Failed to unmarshal message",
			zap.Int("retryCount", retryCount+1),
			zap.Int("maxRetries", maxRetries),
			zap.Error(err))

		msg.Nack(false, true)
		return
	}

	log.Info("Processing order",
		zap.String("orderID", event.OrderID),
		zap.Int("retryCount", retryCount+1),
		zap.Int("maxRetries", maxRetries))

	if err := processOrder(event, log); err != nil {
		log.Error("Failed to process order",
			zap.String("orderID", event.OrderID),
			zap.Int("retryCount", retryCount+1),
			zap.Int("maxRetries", maxRetries),
			zap.Error(err))
		msg.Nack(false, true)
		return
	}

	log.Info("Order processed successfully: ", zap.String("orderID", event.OrderID))
	msg.Ack(false)
}

func processOrder(event events.OrderCreatedEvent, log *zap.Logger) error {
	log.Info("Processing order with total", zap.String("orderID", event.OrderID), zap.Float64("total", event.Total))
	return nil
}

func getRetryCount(msg amqp.Delivery) int {
	if count, ok := msg.Headers["x-retry-count"].(int32); ok {
		return int(count)
	}

	if deaths, ok := msg.Headers["x-death"].([]interface{}); ok && len(deaths) > 0 {
		if death, ok := deaths[0].(amqp.Table); ok {
			if count, ok := death["count"].(int64); ok {
				return int(count)
			}
		}
	}

	return 0
}
