package app

import (
	"github.com/factotum/moneymaker/account-update-service/pkg/accounts"
	"github.com/jaydamon/moneymakergocloak"
)

func (a *App) InitializeRabbitReceivers() {

	goCloakMiddleWare := moneymakergocloak.NewMiddleWare(a.Config.KeyCloakConfig)

	a.RabbitConnection.DeclareExchange("account_update")

	plaidApi := accounts.NewApiService(a.Config.Plaid.Config)
	go a.RabbitConnection.ReceiveMessages(
		"account_refresh",
		accounts.NewReceiver(a.AccountRepository, a.RabbitConnection, goCloakMiddleWare, plaidApi).HandleAccountUpdateEvent)
}
