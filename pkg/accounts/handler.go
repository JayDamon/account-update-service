package accounts

import (
	"encoding/json"
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/config"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/jaydamon/moneymakergocloak"
	"github.com/jaydamon/moneymakerplaid"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type AccountHandler struct {
	rabbitConnection  moneymakerrabbit.Connector
	goCloakMiddleWare moneymakergocloak.Middleware
	plaidHandler      moneymakerplaid.Handler
	config            *config.Config
}

type Handler interface {
	HandleAccountRefreshEvent(msg *amqp091.Delivery)
}

func NewHandler(
	rabbitConnection moneymakerrabbit.Connector,
	goCloakMiddleWare moneymakergocloak.Middleware,
	plaidHandler moneymakerplaid.Handler,
	config *config.Config) Handler {
	return &AccountHandler{
		rabbitConnection:  rabbitConnection,
		goCloakMiddleWare: goCloakMiddleWare,
		plaidHandler:      plaidHandler,
		config:            config,
	}
}

func (handler *AccountHandler) HandleAccountRefreshEvent(msg *amqp091.Delivery) {

	log.Println("Received Message from account-refresh queue")

	err := handler.goCloakMiddleWare.AuthorizeMessage(msg)
	if err != nil {
		fmt.Printf("unauthorized message. %s\n", err)
		// TODO: Send to DLQ
		return
	}
	token, err := moneymakergocloak.GetAuthorizationHeaderFromMessage(msg)
	if err != nil {
		fmt.Printf("error when extracting token from request. %s\n", err)
		// TODO: Send to DLQ
		return
	}
	userId, err := handler.goCloakMiddleWare.ExtractUserIdFromToken(&token)
	if err != nil {
		fmt.Printf("error extracting user id from jwt token. %s\n", err)
		// TODO: Send to DLQ
		return
	}
	fmt.Println("successfully authorized message")

	log.Printf("Processing body %s\n", msg.Body)
	var privateToken users.PrivateToken
	err = json.Unmarshal(msg.Body, &privateToken)
	if err != nil {
		log.Printf("Unable to unmarshal body to Private Token object \n%s\n", msg.Body)
		// TODO: Send to DLQ
		return
	}
	log.Printf("Unmarshalled message body to Private Token object %+v\n", privateToken)

	// TODO: Handle returning no accounts as either nil or empty array
	accounts, err := handler.plaidHandler.GetAccountsForItem(*privateToken.PrivateToken)
	if err != nil {
		log.Printf("Unable to retrieve accounts details \n%s\n", err)
		// TODO: Send to DLQ
		return
	}

	log.Printf("Found accounts. Emitting to Account Update Queue. \n%+v\n", accounts)
	err = emitAccountUpdates(handler.rabbitConnection, accounts, privateToken.ItemId, &token, &userId)
	if err != nil {
		log.Printf("Unable to send all account updates \n%s\n", err)
		// TODO: Send to DLQ
		return
	}
}
