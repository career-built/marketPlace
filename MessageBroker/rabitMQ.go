package MessageBroker

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// AMQP messaging implementation
type RabbitMQBroker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewRabbitMQBroker creates a new instance of RabbitMQBroker.
func NewRabbitMQBroker(connURL string) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(connURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQBroker{
		conn: conn,
		ch:   ch,
	}, nil
}

// -------- Publisher (Sender): ---------
// PublishMessages publishes messages to the specified queue.
func (obj *RabbitMQBroker) PublishMessages(exchange string, queueName string, messages []string) error {
	fmt.Println("PublishMessages Started")
	fmt.Println("PublishMessages initiate Message Tracer if exists")

	time.Sleep(100 * time.Millisecond)
	for _, message := range messages {
		err := obj.ch.Publish(
			"",        // exchange
			queueName, // routing key (queue name)
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------- Consumer (Receiver): ---------
// ConsumeMessages consumes messages from the specified queue and invokes the handler for each message.
func (obj *RabbitMQBroker) ConsumeMessages(queueName string, handler func(message string)) error {
	fmt.Println("ConsumeMessages")

	// Declare the queue
	_, err := obj.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Consume messages
	msgs, err := obj.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack (acknowledgment)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// Handle incoming messages
	go func() {
		for msg := range msgs {
			handler(string(msg.Body))
		}
	}()

	return nil
}

// Close closes the RabbitMQ connection and channel.
func (obj *RabbitMQBroker) Close() error {
	if obj.ch != nil {
		if err := obj.ch.Close(); err != nil {
			return err
		}
	}
	if obj.conn != nil {
		if err := obj.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
