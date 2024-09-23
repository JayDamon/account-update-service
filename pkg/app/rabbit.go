package app

import (
	"github.com/factotum/moneymaker/account-update-service/pkg/accounts"
	"github.com/jaydamon/moneymakergocloak"
)

func (a *App) InitializeRabbitReceivers() {

	goCloakMiddleWare := moneymakergocloak.NewMiddleWare(a.Config.KeyCloakConfig)

	a.RabbitConnection.DeclareExchange("account_update")

	plaidHandler := a.Config.Plaid.NewPlaidHandler()

	go a.RabbitConnection.ReceiveMessages(
		"account_refresh",
		accounts.NewHandler(a.RabbitConnection, goCloakMiddleWare, plaidHandler, a.Config).HandleAccountRefreshEvent)
}
