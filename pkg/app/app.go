package app

import (
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/config"
	"github.com/jaydamon/moneymakerrabbit"
	"log"
	"net/http"
)

type App struct {
	Server           *http.Server
	RabbitConnection moneymakerrabbit.Connector
	Config           *config.Config
}

func NewApplication() *App {
	return &App{
		Config: config.GetConfig(),
	}
}

func (a *App) Initialize() {
	a.Server = &http.Server{
		Addr: fmt.Sprintf(":%s", a.Config.HostPort),
	}
	a.RabbitConnection = a.Config.Rabbit.Connect()
}

func (a *App) Run() {
	appName := a.Config.ApplicationName
	if appName == "" {
		appName = "unnamed service"
	}
	log.Printf("Starting \"%s\" service on port %s\n", appName, a.Config.HostPort)
	err := a.Server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
