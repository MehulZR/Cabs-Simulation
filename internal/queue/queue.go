package queue

import (
	"context"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var NewRideRequestedTopic = "new_ride_requested"
var RideCompletedTopic = "ride_completed"
var DriverAvailableTopic = "driver_available"

var rabbitMQConn *amqp.Connection
var rabbitMQConnErr error
var once sync.Once

type RabbitMQClient struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func CreateNewRabbitMQClient() (*RabbitMQClient, error) {
	once.Do(func() {
		rabbitMQConn, rabbitMQConnErr = amqp.Dial(os.Getenv("RABBITMQ_CONN_STRING"))
	})

	if rabbitMQConnErr != nil {
		return nil, rabbitMQConnErr
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		return nil, err
	}

	rabbitMQClient := RabbitMQClient{Conn: rabbitMQConn, Ch: ch}
	return &rabbitMQClient, nil
}

func (rc RabbitMQClient) Close() error {
	return rc.Ch.Close()
}

func (rc RabbitMQClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.Ch.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	return err
}

func (rc RabbitMQClient) Send(ctx context.Context, queueName string, msg amqp.Publishing) error {
	return rc.Ch.PublishWithContext(ctx, "", queueName, true, false, msg)
}

func (rc RabbitMQClient) Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.Ch.Consume(queueName, consumer, autoAck, false, false, false, nil)
}
