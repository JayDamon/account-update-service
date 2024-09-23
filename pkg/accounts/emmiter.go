package accounts

import (
	"fmt"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/plaid/plaid-go/plaid"
	"log"
)

func emitAccountUpdates(
	rabbitConnection moneymakerrabbit.Connector,
	accountUpdates *plaid.AccountsGetResponse,
	itemId *string,
	bearerToken *string,
	userId *string) error {

	accounts := convertAccountResponseToAccountList(accountUpdates, itemId, userId)

	headers := make(map[string]interface{})
	headers["Authorization"] = *bearerToken

	var errors []error

	for _, a := range *accounts {
		log.Printf("sending message to account_update exchange \n%+v", a)
		err := rabbitConnection.SendMessage(a, headers, "application/json", "", "account_update")
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send one or more messages %v\n", errors)
	}

	return nil
}
