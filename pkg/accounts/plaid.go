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
	GetAccountBalancesForItem(ctx context.Context, accountBalancesGetReq *plaid.AccountsBalanceGetRequest) (plaid.AccountsGetResponse, *http.Response, error)
	GetInstitutionById(ctx context.Context, institutionRequest *plaid.InstitutionsGetByIdRequest) (plaid.InstitutionsGetByIdResponse, *http.Response, error)
}

func (api *PlaidAccountApi) GetAccountsForItem(
	ctx context.Context,
	accountsGetRequest *plaid.AccountsGetRequest,
) (plaid.AccountsGetResponse, *http.Response, error) {
	return api.plaidApi.AccountsGet(ctx).AccountsGetRequest(*accountsGetRequest).Execute()
}

func (api *PlaidAccountApi) GetAccountBalancesForItem(ctx context.Context, accountBalancesGetReq *plaid.AccountsBalanceGetRequest) (plaid.AccountsGetResponse, *http.Response, error) {
	return api.plaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(*accountBalancesGetReq).Execute()
}

func (api *PlaidAccountApi) GetInstitutionById(ctx context.Context, institutionRequest *plaid.InstitutionsGetByIdRequest) (plaid.InstitutionsGetByIdResponse, *http.Response, error) {
	return api.plaidApi.InstitutionsGetById(ctx).InstitutionsGetByIdRequest(*institutionRequest).Execute()
}

func NewApiService(config *plaid.Configuration) ApiService {
	return &PlaidAccountApi{
		plaidApi: *plaid.NewAPIClient(config).PlaidApi,
	}
}
