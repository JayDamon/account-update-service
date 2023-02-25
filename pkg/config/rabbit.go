package config

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Connection struct {
	Connection *amqp091.Connection
}

type MessageHandlerFunc func(msg *amqp091.Delivery)

// ReceiveMessages should be declared as a goroutine to ensure it does not block application startup
func (conn *Connection) ReceiveMessages(queueName string, handler MessageHandlerFunc) {

	ch := openChannel(conn.Connection)
	defer ch.Close()

	_, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			handler(&d)
		}
	}()

	log.Printf(" [*] Waiting for messages from queue %s\n", queueName)
	<-forever
}
