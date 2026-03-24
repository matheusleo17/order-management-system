package messaging

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMqPublisher struct {
	channel *amqp.Channel
}

func NewRabbitMQPublisher() (*RabbitMqPublisher, error) {
	url := os.Getenv("RABBITMQ_URI")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMqPublisher{
		channel: ch,
	}, nil
}

func (p *RabbitMqPublisher) Publish(queue string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	dlqName := queue + ".dlq"

	_, err = p.channel.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": dlqName,
	}

	_, err = p.channel.QueueDeclare(queue, true, false, false, false, args)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return p.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
