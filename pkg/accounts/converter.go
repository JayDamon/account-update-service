package accounts

import (
	"github.com/plaid/plaid-go/plaid"
	"log"
)

func convertAccountResponseToAccountList(
	accountUpdates *plaid.AccountsGetResponse,
	itemId *string,
	userId *string,
	isNew *bool) *[]Account {

	accounts := make([]Account, len(accountUpdates.Accounts))

	au := *accountUpdates

	for i, a := range au.Accounts {
		log.Printf("Converting account\n    %+v", a)
		account := convertAccountBaseToAccount(a)
		account.ItemId = itemId
		account.TenantId = userId
		account.IsNew = isNew
		accounts[i] = account
		log.Printf("Account converted\n    %+v", &accounts[i])
	}

	return &accounts
}

func convertAccountBaseToAccount(accountBase plaid.AccountBase) Account {

	at := string(accountBase.GetType())
	st := string(accountBase.GetSubtype())

	account := Account{}
	account.PlaidAccountId = &accountBase.AccountId
	account.Name = &accountBase.Name
	account.Mask = accountBase.Mask.Get()
	account.OfficialName = accountBase.OfficialName.Get()
	account.AvailableBalance = accountBase.Balances.Available.Get()
	account.CurrentBalance = accountBase.Balances.Current.Get()
	account.Limit = accountBase.Balances.Limit.Get()
	account.AccountTypeName = &at
	if st == "" {
		account.AccountSubTypeName = nil
	} else {
		account.AccountSubTypeName = &st
	}

	return account
}
