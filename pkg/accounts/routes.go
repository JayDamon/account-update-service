package accounts

import (
	"github.com/factotum/moneymaker/account-update-service/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jaydamon/moneymakergocloak"
	"net/http"
)

func CreateRoutes(config *config.Config, controller *Controller, configureCors bool) http.Handler {
	router := chi.NewRouter()

	if configureCors {
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	keycloakMiddleware := moneymakergocloak.NewMiddleWare(config.KeyCloakConfig)
	router.Use(keycloakMiddleware.AuthorizeHttpRequest)
	router.Use(middleware.Heartbeat("/health"))

	AddRoutes(router, controller)

	return router
}

func AddRoutes(mux *chi.Mux, controller *Controller) {
	mux.Get("/v1/accounts", controller.GetALlAccounts)
	mux.Patch("/v1/accounts/{id}/name", controller.UpdateAccountName)
}
