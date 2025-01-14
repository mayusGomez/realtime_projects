package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQClient struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRabbitClient(connStr string, queues []string) (*RabbitMQClient, error) {
	conn, err := amqp091.Dial(connStr)
	if err != nil {
		log.Printf("failed to connect to RabbitMQ server: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("failed to open a channel: %v", err)
		return nil, err
	}

	err = declareQueues(ch, queues)
	if err != nil {
		log.Printf("failed to declare queues: %v", err)
		return nil, err
	}

	return &RabbitMQClient{ch: ch, conn: conn}, nil
}

func (rc *RabbitMQClient) Close() {
	err := rc.ch.Close()
	if err != nil {
		log.Printf("failed to close channel: %v", err)
	}

	err = rc.conn.Close()
	if err != nil {
		log.Printf("failed to close connection: %v", err)
	}
}

func (rc *RabbitMQClient) SubscribeToQueue(queue string) (<-chan amqp091.Delivery, error) {
	msgs, err := rc.ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func declareQueues(ch *amqp091.Channel, queues []string) error {
	for _, queue := range queues {
		log.Printf("declaring queue %q", queue)
		_, err := ch.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("failed to declare a queue: %v", err)
			return err
		}
	}

	return nil
}

func (rc *RabbitMQClient) Publish(queue string, body []byte) error {
	err := rc.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("failed to publish a message: %v", err)
		return err
	}

	return nil
}
