package main

import (
	"github.com/factotum/moneymaker/account-update-service/pkg/app"
	"log"
)

func main() {
	application := app.NewApplication()

	log.Printf("Initializing application\n")

	application.Initialize()

	defer application.RabbitConnection.Connection.Close()

	application.InitializeRabbitReceivers()

	application.Run()

}