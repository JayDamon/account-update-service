package accounts

import (
	"fmt"
	tools "github.com/jaydamon/http-toolbox"
	"github.com/jaydamon/moneymakergocloak"
	"net/http"
)

type Controller struct {
	AccountRepository Repository
	keycloakConfig    *moneymakergocloak.Configuration
}

func NewController(accountRepository Repository, keycloakConfig *moneymakergocloak.Configuration) *Controller {
	return &Controller{
		AccountRepository: accountRepository,
		keycloakConfig:    keycloakConfig,
	}
}

func (controller *Controller) GetALlAccounts(w http.ResponseWriter, r *http.Request) {
	userId, err := moneymakergocloak.ExtractUserIdFromRequest(r, controller.keycloakConfig)
	if err != nil {
		fmt.Println("Error extracting bearer userId from request")
		tools.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	aa, err := controller.AccountRepository.GetAccountsForUser(userId)
	if err != nil {
		fmt.Println("error retrieving all transactions", err)
		tools.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tools.Respond(w, http.StatusOK, aa)
}

func (controller *Controller) UpdateAccountName(w http.ResponseWriter, r *http.Request) {

	//a, err := controller.AccountRepository.UP
	
}
