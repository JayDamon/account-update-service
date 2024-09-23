package accounts

import (
	"github.com/plaid/plaid-go/plaid"
)

func convertAccountResponseToAccountList(accountUpdates *plaid.AccountsGetResponse, itemId *string, userId *string) *[]Account {

	accounts := make([]Account, len(accountUpdates.Accounts))

	for i, a := range accountUpdates.Accounts {
		account := convertAccountBaseToAccount(&a)
		account.ItemId = itemId
		account.TenantId = userId
		accounts[i] = *account
	}

	return &accounts
}

func convertAccountBaseToAccount(accountBase *plaid.AccountBase) *Account {

	accType := &accountBase.Type

	account := Account{}
	account.AccountId = &accountBase.AccountId
	account.Name = &accountBase.Name
	account.OfficialName = accountBase.OfficialName.Get()
	account.AvailableBalance = accountBase.Balances.Available.Get()
	account.CurrentBalance = accountBase.Balances.Current.Get()
	account.Limit = accountBase.Balances.Limit.Get()
	account.AccountType = accType
	account.AccountSubType = accountBase.Subtype.Get()

	return &account
}
