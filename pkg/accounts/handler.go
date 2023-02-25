package accounts

import (
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/config"
	"github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	rabbitConnection *config.Connection
}

func NewHandler(rabbitConnection *config.Connection) *Handler {
	return &Handler{
		rabbitConnection: rabbitConnection,
	}
}

func (handler *Handler) HandleAccountRefreshEvent(msg *amqp091.Delivery) {
	fmt.Printf("Message Received: %s\n", msg.Body)
}
