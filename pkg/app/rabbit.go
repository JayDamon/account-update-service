package app

import (
	"github.com/factotum/moneymaker/account-update-service/pkg/accounts"
	"github.com/jaydamon/moneymakergocloak"
)

func (a *App) InitializeRabbitReceivers() {

	goCloakMiddleWare := moneymakergocloak.NewMiddleWare(a.Config.KeyCloakConfig)

	go a.RabbitConnection.ReceiveMessages(
		"account_refresh",
		accounts.NewHandler(a.RabbitConnection, goCloakMiddleWare, a.Config).HandleAccountRefreshEvent)
}
