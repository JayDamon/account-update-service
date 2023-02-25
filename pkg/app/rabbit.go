package app

import (
	"github.com/factotum/moneymaker/account-update-service/pkg/accounts"
)

func (a *App) InitializeRabbitReceivers() {
	go a.RabbitConnection.ReceiveMessages(
		"account_refresh",
		accounts.NewHandler(a.RabbitConnection).HandleAccountRefreshEvent)
}
