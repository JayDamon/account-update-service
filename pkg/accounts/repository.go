package accounts

type Repository interface {
	GetAccountsForUser(tenantId string) (*[]Account, error)
	InsertNewAccounts(ai *AccountItem) (*AccountItem, error)
	UpdateTransactionName(a *Account) (*Account, error)
}
