package main

import (
	"context"
	"github.com/factotum/moneymaker/account-update-service/pkg/app"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	application := app.NewApplication()

	log.Printf("Initializing application\n")

	application.Initialize()

	defer application.RabbitConnection.Connection.Close()

	application.InitializeRabbitReceivers()

	application.Run()

}

func sendMessage(ch *amqp091.Channel) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err := ch.PublishWithContext(ctx,
		"",      // exchange
		"hello", // routing key
		false,   // mandatory
		false,   // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
