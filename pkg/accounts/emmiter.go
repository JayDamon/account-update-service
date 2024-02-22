package accounts

import (
	"github.com/plaid/plaid-go/plaid"
)

func (handler *Handler) emitAccountUpdates(accountUpdates *plaid.AccountsGetResponse, bearerToken *string) error {

	// Convert to account list

	// Send to message queue using rabbitmq on EmmitAccount event

	return nil
}

