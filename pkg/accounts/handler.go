package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/jaydamon/moneymakergocloak"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type AccountHandler struct {
	rabbitConnection  moneymakerrabbit.Connector
	goCloakMiddleWare moneymakergocloak.Middleware
	plaidApi          ApiService
}

type Handler interface {
	HandleAccountUpdateEvent(msg *amqp091.Delivery)
}

func NewHandler(
	rabbitConnection moneymakerrabbit.Connector,
	goCloakMiddleWare moneymakergocloak.Middleware,
	plaidApi ApiService) Handler {
	return &AccountHandler{
		rabbitConnection:  rabbitConnection,
		goCloakMiddleWare: goCloakMiddleWare,
		plaidApi:          plaidApi,
	}
}

func (handler *AccountHandler) HandleAccountUpdateEvent(msg *amqp091.Delivery) {

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

	if userId != *privateToken.UserId {
		log.Printf("invalid private token. user id does not match oauth token")
		// TODO: Send to DLQ
		return
	}
	ctx := context.Background()
	accountsGetRequest := *plaid.NewAccountsGetRequest(*privateToken.PrivateToken)
	accounts, _, err := handler.plaidApi.GetAccountsForItem(ctx, &accountsGetRequest)
	if err != nil {
		log.Printf("Unable to retrieve accounts details \n%s\n", err)
		// TODO: Send to DLQ
		return
	}
	// TODO: If account is existing, get balances for them

	log.Printf("Found accounts. Emitting to Account Update Queue. \n%+v\n", accounts)
	err = emitAccountUpdates(handler.rabbitConnection, &accounts, privateToken.Cursor, &token, &privateToken)
	if err != nil {
		log.Printf("Unable to send all account updates \n%s\n", err)
		// TODO: Send to DLQ
		return
	}
}
