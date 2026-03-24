package messaging

import "github.com/streadway/amqp"

type RabbitMQConsumer struct {
	channel *amqp.Channel
}

func NewRabbitMQConsumer(url string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		channel: ch,
	}, nil

}
func (c *RabbitMQConsumer) Consume(queue string) (<-chan amqp.Delivery, error) {
	_, err := c.channel.QueueDeclare(
		"order.created.dlq",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": "order.created.dlq",
	}
	_, err = c.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	return msgs, err
}
