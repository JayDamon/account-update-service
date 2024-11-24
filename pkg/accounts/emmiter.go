package accounts

import (
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/plaid/plaid-go/plaid"
	"log"
)

func emitAccountUpdates(
	rabbitConnection moneymakerrabbit.Connector,
	accountUpdates *plaid.AccountsGetResponse,
	cursor *string,
	bearerToken *string,
	privateToken *users.PrivateToken) error {

	accounts := convertAccountResponseToAccountList(
		accountUpdates,
		privateToken.ItemId,
		privateToken.UserId,
		privateToken.IsNew)

	ai := AccountItem{
		ItemId:   privateToken.ItemId,
		Cursor:   cursor,
		Accounts: accounts,
	}

	headers := make(map[string]interface{})
	headers["Authorization"] = *bearerToken
	headers["PlaidToken"] = *privateToken.PrivateToken

	log.Printf("sending message to account_update exchange \n%+v", ai)
	err := rabbitConnection.SendMessage(ai, headers, "application/json", "", "account_update")
	if err != nil {
		return fmt.Errorf("failed to send messages %v\n", err)
	}

	return nil
}
