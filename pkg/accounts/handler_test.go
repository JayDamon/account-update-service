package accounts

import (
	"github.com/plaid/plaid-go/plaid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestHandleAccountRefreshEvent_SuccessPath(t *testing.T) {

	conn := &TestConnector{}
	mid := &TestMiddleware{}
	mid.extractUserIdFromTokenReturn = "userId"
	plaidHandler := &TestApiService{}
	testRepository := &TestRepository{}
	//plaidHandler.response = createTestAccountsGetResponse()
	handler := createTestHandler(conn, mid, plaidHandler, testRepository)

	msg := createTestMessage()

	err := handler.HandleAccountUpdateEvent(msg)

	assert.Nil(t, err)
	assert.Equal(t, 1, mid.authorizeMessagesCount)
	assert.Equal(t, 0, mid.authorizeHttpRequestCount)
	assert.Equal(t, 1, mid.extractUserIdFromTokenCount)
	assert.Equal(t, 1, plaidHandler.getAccountsForItemCount)
	assert.Equal(t, "token", plaidHandler.request.GetAccessToken())
	assert.Equal(t, 1, conn.sendMessageCount)
	assert.Equal(t, "testToken", conn.headers["Authorization"])
	assert.Equal(t, "token", conn.headers["PlaidToken"])
}

func createTestMessage() *amqp091.Delivery {

	headers := make(map[string]interface{})
	headers["Authorization"] = "testToken"

	body := []byte("{\"id\":\"userId\",\"privateToken\":\"token\",\"itemId\":\"itemId\",\"isNew\":true,\"cursor\":\"testCursor\"}")

	return &amqp091.Delivery{
		nil,
		headers,
		"application/json",
		"utf-8",
		1,
		0,
		"",
		"",
		"",
		"",
		time.Now(),
		"",
		"",
		"",
		"",
		1,
		1,
		false,
		"account_refresh",
		"",
		body,
	}

}

type TestMiddleware struct {
	authorizeMessagesCount       int
	authorizeHttpRequestCount    int
	extractUserIdFromTokenCount  int
	extractUserIdFromTokenReturn string
	extractUserIdFromTokenError  error
}

func (middleware *TestMiddleware) AuthorizeMessage(msg *amqp091.Delivery) error {
	middleware.authorizeMessagesCount++
	return nil
}

func (middleware *TestMiddleware) AuthorizeHttpRequest(request http.Handler) http.Handler {
	middleware.authorizeHttpRequestCount++
	return nil
}

func (middlware *TestMiddleware) ExtractUserIdFromToken(token *string) (string, error) {
	middlware.extractUserIdFromTokenCount++
	return middlware.extractUserIdFromTokenReturn, middlware.extractUserIdFromTokenError
}

type TestPlaidHandler struct {
	getAccountsForItemCount int
	privateToken            string
	response                *plaid.AccountsGetResponse
}

func (handler *TestPlaidHandler) GetAccountsForItem(privateToken string) (*plaid.AccountsGetResponse, error) {
	handler.getAccountsForItemCount++
	handler.privateToken = privateToken
	return handler.response, nil
}

type TestRepository struct {
}

func (r *TestRepository) GetAccountsForUser(tenantId string) (*[]Account, error) {
	return nil, nil
}

func (r *TestRepository) InsertNewAccounts(ai *AccountItem) (*AccountItem, error) {
	return nil, nil
}

func (r *TestRepository) UpdateTransactionName(a *Account) (*Account, error) {
	return nil, nil
}
