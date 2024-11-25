package accounts

import (
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/jaydamon/moneymakerrabbit"
	"log"
)

func emitAccountUpdates(rabbitConnection moneymakerrabbit.Connector, ai *AccountItem, bearerToken *string, privateToken *users.PrivateToken) error {

	headers := make(map[string]interface{})
	headers["Authorization"] = *bearerToken
	headers["PlaidToken"] = *privateToken.PrivateToken

	log.Printf("sending message to account_update exchange \n    %+v\n    accounts: %+v\n", ai, ai.Accounts)
	err := rabbitConnection.SendMessage(ai, headers, "application/json", "", "account_update")
	if err != nil {
		return fmt.Errorf("failed to send messages %v\n", err)
	}

	return nil
}
