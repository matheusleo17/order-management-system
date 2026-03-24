package messaging

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type RabbitMqPublisher struct {
	channel *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*RabbitMqPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}
	return &RabbitMqPublisher{
		channel: ch,
	}, nil

}

func (p *RabbitMqPublisher) Publish(queue string, message interface{}) error {

	body, err := json.Marshal(message)

	if err != nil {
		return err
	}
	dlqName := queue + ".dlq"

	_, err = p.channel.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		nil,
	)

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": dlqName,
	}

	_, err = p.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		args,
	)
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
