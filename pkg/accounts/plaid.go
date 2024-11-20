package accounts

import (
	"context"
	"github.com/plaid/plaid-go/plaid"
	"net/http"
)

type PlaidAccountApi struct {
	plaidApi plaid.PlaidApiService
}

type ApiService interface {
	GetAccountsForItem(ctx context.Context, accountsGetRequest *plaid.AccountsGetRequest) (plaid.AccountsGetResponse, *http.Response, error)
}

func (api *PlaidAccountApi) GetAccountsForItem(
	ctx context.Context,
	accountsGetRequest *plaid.AccountsGetRequest,
) (plaid.AccountsGetResponse, *http.Response, error) {
	return api.plaidApi.AccountsGet(ctx).AccountsGetRequest(*accountsGetRequest).Execute()
}

func NewApiService(config *plaid.Configuration) ApiService {
	return &PlaidAccountApi{
		plaidApi: *plaid.NewAPIClient(config).PlaidApi,
	}
}
