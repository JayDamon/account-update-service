package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/jaydamon/moneymakergocloak"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type AccountReceiver struct {
	accountRepository Repository
	rabbitConnection  moneymakerrabbit.Connector
	goCloakMiddleWare moneymakergocloak.Middleware
	plaidApi          ApiService
}

type Receiver interface {
	HandleAccountUpdateEvent(msg *amqp091.Delivery) error
}

func NewReceiver(
	accountRepository Repository,
	rabbitConnection moneymakerrabbit.Connector,
	goCloakMiddleWare moneymakergocloak.Middleware,
	plaidApi ApiService) Receiver {

	return &AccountReceiver{
		accountRepository: accountRepository,
		rabbitConnection:  rabbitConnection,
		goCloakMiddleWare: goCloakMiddleWare,
		plaidApi:          plaidApi,
	}
}

func (handler *AccountReceiver) HandleAccountUpdateEvent(msg *amqp091.Delivery) error {

	log.Println("Received Message from account-refresh queue")

	err := handler.goCloakMiddleWare.AuthorizeMessage(msg)
	if err != nil {
		fmt.Printf("unauthorized message. %s\n", err)
		return err
	}
	token, err := moneymakergocloak.GetAuthorizationHeaderFromMessage(msg)
	if err != nil {
		fmt.Printf("error when extracting token from request. %s\n", err)
		return err
	}
	userId, err := handler.goCloakMiddleWare.ExtractUserIdFromToken(&token)
	if err != nil {
		fmt.Printf("error extracting user id from jwt token. %s\n", err)
		return err
	}
	fmt.Println("successfully authorized message")

	log.Printf("Processing body %s\n", msg.Body)
	var privateToken users.PrivateToken
	err = json.Unmarshal(msg.Body, &privateToken)
	if err != nil {
		log.Printf("Unable to unmarshal body to Private Token object \n%s\n", msg.Body)
		return err
	}
	log.Printf("Unmarshalled message body to Private Token object %+v\n", privateToken)

	if userId != *privateToken.UserId {
		log.Printf("invalid private token. user id does not match oauth token")
		return err
	}
	ctx := context.Background()
	accountsGetRequest := *plaid.NewAccountsGetRequest(*privateToken.PrivateToken)
	accountsResponse, _, err := handler.plaidApi.GetAccountsForItem(ctx, &accountsGetRequest)
	if err != nil {
		log.Printf("Unable to retrieve accounts details \n%s\n", err)
		return err
	}

	if !*privateToken.IsNew {
		balancesGetReq := plaid.NewAccountsBalanceGetRequest(*privateToken.PrivateToken)
		balancesGetResp, _, err := handler.plaidApi.GetAccountBalancesForItem(ctx, balancesGetReq)
		if err != nil {
			log.Printf("Unable to retrieve account balances \n%s\n", err)
		}
		var bb map[string]plaid.AccountBase
		for _, b := range balancesGetResp.Accounts {
			bb[b.AccountId] = b
		}
		for _, a := range accountsResponse.Accounts {
			a.Balances = bb[a.AccountId].Balances
		}
	}

	log.Printf("Found accounts. Emitting to Account Update Queue. \n%+v\n", accountsResponse)
	accounts := convertAccountResponseToAccountList(
		&accountsResponse,
		privateToken.ItemId,
		privateToken.UserId,
		privateToken.IsNew)

	ai := AccountItem{
		ItemId:   privateToken.ItemId,
		TenantId: privateToken.UserId,
		Cursor:   privateToken.Cursor,
		Accounts: accounts,
	}

	if accountsResponse.Item.GetInstitutionId() != "" {
		cc := []plaid.CountryCode{plaid.COUNTRYCODE_US}
		institutionRequest := plaid.NewInstitutionsGetByIdRequest(accountsResponse.Item.GetInstitutionId(), cc)
		iResp, _, err := handler.plaidApi.GetInstitutionById(ctx, institutionRequest)
		if err != nil {
			log.Printf("Unable to retrieve institution details \n%s\n", err)
			return err
		}
		institution := iResp.Institution
		url := institution.GetUrl()
		pc := institution.GetPrimaryColor()
		l := institution.GetLogo()

		ai.InstitutionId = &institution.InstitutionId
		ai.InstitutionName = &institution.Name
		ai.Url = &url
		ai.PrimaryColor = &pc
		ai.Logo = &l
	}

	_, err = handler.accountRepository.InsertNewAccounts(&ai)
	if err != nil {
		log.Printf("Unable to save account details \n%s\n", err)
		return err
	}

	err = emitAccountUpdates(handler.rabbitConnection, &ai, &token, &privateToken)
	if err != nil {
		log.Printf("Unable to send all account updates \n%s\n", err)
		return err
	}
	return nil
}
