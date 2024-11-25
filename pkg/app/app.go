package app

import (
	"database/sql"
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/accounts"
	"github.com/factotum/moneymaker/account-update-service/pkg/config"
	"github.com/jaydamon/moneymakerrabbit"
	"log"
	"net/http"
)

type App struct {
	Server            *http.Server
	DB                *sql.DB
	AccountRepository accounts.Repository
	AccountController *accounts.Controller
	RabbitConnection  moneymakerrabbit.Connector
	Config            *config.Config
}

func NewApplication() *App {
	return &App{
		Config: config.GetConfig(),
	}
}

func (a *App) Initialize() {
	a.DB = connectToDB(a.Config)
	a.AccountRepository = accounts.NewRepository(a.DB)
	a.AccountController = accounts.NewController(a.AccountRepository, a.Config.KeyCloakConfig)
	a.Server = &http.Server{
		Addr:    fmt.Sprintf(":%s", a.Config.HostPort),
		Handler: accounts.CreateRoutes(a.Config, a.AccountController, a.Config.ConfigureCors),
	}
	performDbMigration(a.DB, a.Config)
	a.RabbitConnection = a.Config.Rabbit.Connect()
}

func (a *App) Run() {

	defer a.DB.Close()

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
