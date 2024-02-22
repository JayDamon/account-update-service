package accounts

import (
	"context"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/plaid/plaid-go/plaid"
	"log"
)

func (handler *Handler) GetAccountsForItem(token *users.PrivateToken) (*plaid.AuthGetResponse, error) {

	ctx := context.Background()
	plaidConfig := handler.config.Plaid
	plaidClient := plaidConfig.Client

	authGetRequest := *plaid.NewAuthGetRequest(*token.PrivateToken)

	authGetResp, _, err := plaidClient.PlaidApi.AuthGet(ctx).AuthGetRequest(
		authGetRequest,
	).Execute()
	if err != nil {
		log.Printf("Unable to get account details \n%+v\n", err)
		return nil, err
	}
	log.Printf("Retrieved Auth Response \n%+v\n", authGetResp)

	return &authGetResp, nil
}