package commons

import (
	"fmt"
	"log"

	rabbitmq "github.com/wagslane/go-rabbitmq"
)

type rabbitMQ struct {
	address   string
	username  string
	password  string
	Conn      *rabbitmq.Conn
	Publisher *rabbitmq.Publisher
}

func NewRabbitMQConnection(address, username, password string) *rabbitMQ {
	return &rabbitMQ{
		address:  address,
		username: username,
		password: password,
	}
}

func (r *rabbitMQ) Connect() error {
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s", r.username, r.password, r.address),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return err
	}

	r.Conn = conn

	return nil
}

func (r *rabbitMQ) CloseConnection(conn *rabbitmq.Conn) error {
	return conn.Close()
}

func (r *rabbitMQ) NewRabbitMQPublisher() error {
	// establish connection
	err := r.Connect()
	if err != nil {
		return err
	}
	// defer r.CloseConnection(conn)

	publisher, err := rabbitmq.NewPublisher(
		r.Conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}
	// defer publisher.Close()

	publisher.NotifyReturn(func(r rabbitmq.Return) {
		log.Printf("message returned from server: %s", string(r.Body))
	})
	publisher.NotifyPublish(func(c rabbitmq.Confirmation) {
		log.Printf("message confirmed from server. tag: %v, ack: %v", c.DeliveryTag, c.Ack)
	})

	r.Publisher = publisher

	return nil
}
