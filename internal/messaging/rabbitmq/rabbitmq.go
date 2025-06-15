package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MesssageQueue interface {
	Publish(exchange, routingKey string, message []byte) error
	Consume(queue string) (<-chan amqp.Delivery, error)
	DeclareExchange(name, kind string) error
	DeclareQueue(name string) (amqp.Queue, error)
	BindQueue(queue, exchange, routingKey string) error
	CLose() error
}

type rabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (MesssageQueue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %v", err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		defer conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err.Error())
	}

	return &rabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *rabbitMQ) Publish(exchange, routingKey string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.channel.PublishWithContext(ctx, exchange, routingKey,
		false, //mandadory
		false, //immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
			Timestamp:   time.Now(),
		},
	)
}

func (r *rabbitMQ) Consume(queue string) (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (r *rabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(name, kind,
		true,  //durable
		false, //auto-deleted
		false, //internal
		false, //no-wait
		nil,   //arguments
	)
}

func (r *rabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(name,
		true,  //durable
		false, //auto-deleted
		false, //exclusive
		false, //no-wait
		nil,   //arguments
	)
}

func (r *rabbitMQ) BindQueue(queue, exchange, routingKey string) error {
	return r.channel.QueueBind(queue, routingKey, exchange,
		false, //no-wait
		nil,   //args
	)
}

func (r *rabbitMQ) CLose() error {
	return nil
}
