package accounts

import "github.com/plaid/plaid-go/plaid"

type Account struct {
	TenantId         *string               `json:"tenantId"`
	Name             *string               `json:"name"`
	AccountId        *string               `json:"id"`
	ItemId           *string               `json:"itemId"`
	OfficialName     *string               `json:"officialName"`
	AvailableBalance *float32              `json:"availableBalance"`
	CurrentBalance   *float32              `json:"currentBalance"`
	Limit            *float32              `json:"limit,omitempty"`
	AccountType      *plaid.AccountType    `json:"accountType"`
	AccountSubType   *plaid.AccountSubtype `json:"accountSubType"`
}

//type Balances struct {
//	Available *float32 `json:"available"`
//	Current *float32 `json:"current"`
//	Limit *float32 `json:"limit"`
//	IsoCurrencyCode *string `json:"iso_currency_code"`
//	UnofficialCurrencyCode *string `json:"unofficial_currency_code"`
//}
