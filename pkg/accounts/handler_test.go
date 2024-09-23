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
	plaidHandler := &TestPlaidHandler{}
	plaidHandler.response = createTestAccountsGetResponse()
	handler := createTestHandler(conn, mid, plaidHandler)

	msg := createTestMessage()

	handler.HandleAccountRefreshEvent(msg)

	assert.Equal(t, 1, mid.authorizeMessagesCount)
	assert.Equal(t, 0, mid.authorizeHttpRequestCount)
	assert.Equal(t, 1, mid.extractUserIdFromTokenCount)
	assert.Equal(t, 1, plaidHandler.getAccountsForItemCount)
	assert.Equal(t, "token", plaidHandler.privateToken)
	assert.Equal(t, 2, conn.sendMessageCount)
	assert.Equal(t, "testToken", conn.headers[0]["Authorization"])
}

func createTestMessage() *amqp091.Delivery {

	headers := make(map[string]interface{})
	headers["Authorization"] = "testToken"

	body := []byte("{\"id\":\"userId\",\"privateToken\":\"token\",\"itemId\":\"itemId\"}")

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
