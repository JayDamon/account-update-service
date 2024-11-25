package accounts

import (
	"context"
	"fmt"
	"github.com/factotum/moneymaker/account-update-service/pkg/users"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEmitAccountUpdates_SuccessPath(t *testing.T) {

	conn := &TestConnector{}

	accountsGetResponse := createTestAccountsGetResponse()
	token := "token"

	itemId := "itemId"
	tenantId := "tenantId"
	isNew := true
	plaidToken := "privateToken"

	privateToken := &users.PrivateToken{
		ItemId:       &itemId,
		PrivateToken: &plaidToken,
		UserId:       &tenantId,
		IsNew:        &isNew,
	}

	aa := convertAccountResponseToAccountList(accountsGetResponse, &itemId, &tenantId, &isNew)

	ais := AccountItem{
		ItemId:   &itemId,
		TenantId: &tenantId,
		Accounts: aa,
	}

	err := emitAccountUpdates(conn, &ais, &token, privateToken)

	assert.Nil(t, err)

	assert.Equal(t, 1, conn.sendMessageCount)
	assert.Equal(t, token, conn.headers["Authorization"])
	assert.Equal(t, "", conn.queue)
	assert.Equal(t, "account_update", conn.exchange)
	assert.Equal(t, "application/json", conn.contentType)

	ai, ok := conn.body.(*AccountItem)
	assert.True(t, ok)
	assert.Equal(t, "itemId", *ai.ItemId)
	assert.Equal(t, 2, len(*ai.Accounts))

	var accountsChecked int
	for _, a := range *ai.Accounts {

		assert.Equal(t, "itemId", *a.ItemId)
		assert.Equal(t, true, *a.IsNew)

		accountsChecked++
	}

	assert.Equal(t, 2, accountsChecked)

}

func TestEmitAccountUpdates_SendMessageFails(t *testing.T) {
	conn := &TestConnector{}
	conn.failTimeCalled = 1

	accountsGetResponse := createTestAccountsGetResponse()
	token := "token"
	itemId := "itemId"
	tenantId := "tenantId"
	isNew := false
	plaidToken := "privateToken"

	privateToken := &users.PrivateToken{
		ItemId:       &itemId,
		PrivateToken: &plaidToken,
		UserId:       &tenantId,
		IsNew:        &isNew,
	}

	aa := convertAccountResponseToAccountList(accountsGetResponse, &itemId, &tenantId, &isNew)

	ais := AccountItem{
		ItemId:   &itemId,
		TenantId: &tenantId,
		Accounts: aa,
	}

	err := emitAccountUpdates(conn, &ais, &token, privateToken)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to send messages failing for 1 test\n", err.Error())

	assert.Equal(t, 1, conn.sendMessageCount)

	ai, ok := conn.body.(*AccountItem)
	assert.True(t, ok)
	assert.Equal(t, "itemId", *ai.ItemId)
	assert.Equal(t, 2, len(*ai.Accounts))

	var accountsChecked int
	for _, a := range *ai.Accounts {

		assert.Equal(t, "itemId", *a.ItemId)

		accountsChecked++
	}

	assert.Equal(t, 2, accountsChecked)
}

func createTestHandler(conn *TestConnector, middleware *TestMiddleware, plaidHandler *TestApiService, testRepository *TestRepository) Receiver {

	handler := NewReceiver(testRepository, conn, middleware, plaidHandler)

	return handler
}

type TestConnector struct {
	body                 interface{}
	headers              map[string]interface{}
	contentType          string
	queue                string
	exchange             string
	receiveMessagesCount int
	sendMessageCount     int
	failTimeCalled       int
}

func (conn *TestConnector) ReceiveMessages(queueName string, handler moneymakerrabbit.MessageHandlerFunc) {
	conn.receiveMessagesCount++
}

func (conn *TestConnector) SendMessage(body interface{}, headers map[string]interface{}, contentType string, queue string, exchange string) error {
	conn.body = body
	conn.headers = headers
	conn.contentType = contentType
	conn.queue = queue
	conn.exchange = exchange
	conn.sendMessageCount++

	if conn.sendMessageCount == conn.failTimeCalled {
		return fmt.Errorf("failing for %o test", conn.sendMessageCount)
	}

	return nil
}

func (conn *TestConnector) Close() {}

func (conn *TestConnector) DeclareExchange(exchangeName string) {}

func (conn *TestConnector) DeclareQueue(queueName string) *amqp091.Queue {
	return nil
}

func (conn *TestConnector) ReceiveMessagesFromExchange(exchangeName string, consumingQueueName string, handler moneymakerrabbit.MessageHandlerFunc) {
}

type TestApiService struct {
	request                 plaid.AccountsGetRequest
	getAccountsForItemCount int
}

func (api *TestApiService) GetAccountsForItem(ctx context.Context, accountsGetRequest *plaid.AccountsGetRequest) (plaid.AccountsGetResponse, *http.Response, error) {
	response := plaid.AccountsGetResponse{}
	api.request = *accountsGetRequest
	api.getAccountsForItemCount++
	return response, nil, nil
}

func (api *TestApiService) GetAccountBalancesForItem(ctx context.Context, accountBalancesGetReq *plaid.AccountsBalanceGetRequest) (plaid.AccountsGetResponse, *http.Response, error) {
	response := plaid.AccountsGetResponse{}
	return response, nil, nil
}

func (api *TestApiService) GetInstitutionById(ctx context.Context, institutionRequest *plaid.InstitutionsGetByIdRequest) (plaid.InstitutionsGetByIdResponse, *http.Response, error) {
	response := plaid.InstitutionsGetByIdResponse{}
	return response, nil, nil
}
