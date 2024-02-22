package accounts

import (
	"encoding/json"
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/config"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/jaydamon/moneymakergocloak"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Handler struct {
	rabbitConnection *moneymakerrabbit.Connection
	goCloakMiddleWare *moneymakergocloak.Middleware
	config *config.Config
}

func NewHandler(
	rabbitConnection *moneymakerrabbit.Connection,
	goCloakMiddleWare *moneymakergocloak.Middleware,
	config *config.Config) *Handler {
	return &Handler{
		rabbitConnection: rabbitConnection,
		goCloakMiddleWare: goCloakMiddleWare,
		config: config,
	}
}

func (handler *Handler) HandleAccountRefreshEvent(msg *amqp091.Delivery) {
	err := handler.goCloakMiddleWare.AuthorizeMessage(msg)
	if err != nil {
		fmt.Printf("unauthorized message. %s\n", err)
		// TODO: Send to DLQ
		return
	}

	fmt.Print("successfully authorized message")

	//bearerToken, err := moneymakergocloak.GetAuthorizationHeaderFromMessage(msg)
	//if err != nil {
	//	log.Printf("Authorization Header not provided \n%s\n", err)
	//	// TODO: Send to DLQ
	//	return
	//}

	log.Printf("Processing body %s", msg.Body)
	var privateToken users.PrivateToken
	err = json.Unmarshal(msg.Body, &privateToken)
	if err != nil {
		log.Printf("Unable to unmarshal body to Private Token object \n%s\n", msg.Body)
		// TODO: Send to DLQ
		return
	}
	log.Printf("Unmarshalled message body to Private Token object %+v", privateToken)

	accounts, err := handler.GetAccountsForItem(&privateToken)
	if err != nil {
		log.Printf("Unable to retrieve accounts details \n%s\n", err)
		// TODO: Send to DLQ
		return
	}

	log.Printf("Found accounts. Emitting to Account Update Queue. \n%+v\n", accounts)

}
