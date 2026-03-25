package messaging

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQConsumer struct {
	channel *amqp.Channel
	log     *zap.Logger
}

func NewRabbitMQConsumer(log *zap.Logger) (*RabbitMQConsumer, error) {
	url := os.Getenv("RABBITMQ_URI")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}
	log.Info("Connecting to RabbitMQ (consumer)")

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	log.Info("RabbitMQ consumer connected")

	return &RabbitMQConsumer{
		channel: ch,
		log:     log,
	}, nil
}

func (c *RabbitMQConsumer) Consume(queue string) (<-chan amqp.Delivery, error) {
	dlqName := queue + ".dlq"

	_, err := c.channel.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare DLQ: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": dlqName,
	}

	_, err = c.channel.QueueDeclare(queue, true, false, false, false, args)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := c.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming: %w", err)
	}

	return msgs, nil
}
