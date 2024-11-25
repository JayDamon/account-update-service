package accounts

type Account struct {
	Id                     *string  `json:"accountId,omitempty"`
	TenantId               *string  `json:"tenantId"`
	FriendlyName           *string  `json:"friendlyName,omitempty"`
	Name                   *string  `json:"name"`
	PlaidAccountId         *string  `json:"plaidAccountId,omitempty"`
	Mask                   *string  `json:"mask,omitempty"`
	ItemId                 *string  `json:"itemId,omitempty"`
	OfficialName           *string  `json:"officialName,omitempty"`
	AvailableBalance       *float32 `json:"availableBalance,omitempty"`
	CurrentBalance         *float32 `json:"currentBalance,omitempty"`
	Limit                  *float32 `json:"limit,omitempty"`
	OfficialCurrencyCode   *string  `json:"officialCurrencyCode,omitempty"`
	UnofficialCurrencyCode *string  `json:"unofficialCurrencyCode,omitempty"`
	IsPrimaryAccount       *bool    `json:"isPrimaryAccount,omitempty"`
	IsInCashFlow           *bool    `json:"isInCashFlow,omitempty"`
	AccountTypeName        *string  `json:"accountType,omitempty"`
	AccountSubTypeName     *string  `json:"accountSubType,omitempty"`
	InstitutionId          *string  `json:"institutionId,omitempty"`
	InstitutionName        *string  `json:"institutionName,omitempty"`
	IsNew                  *bool    `json:"isNew,omitempty"`
}

type AccountItem struct {
	ItemId          *string    `json:"itemId"`
	Cursor          *string    `json:"cursor,omitempty"`
	InstitutionId   *string    `json:"institutionId,omitempty"`
	InstitutionName *string    `json:"institutionName,omitempty"`
	Url             *string    `json:"url,omitempty"`
	PrimaryColor    *string    `json:"primaryColor,omitempty"`
	Logo            *string    `json:"logo,omitempty"`
	TenantId        *string    `json:"tenantId"`
	Accounts        *[]Account `json:"accounts,omitempty"`
}
