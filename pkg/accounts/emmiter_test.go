package accounts

import (
	"fmt"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmitAccountUpdates_SuccessPath(t *testing.T) {

	conn := &TestConnector{}

	accountsGetResponse := createTestAccountsGetResponse()
	token := "token"
	itemId := "itemId"
	tenantId := "tenantId"

	err := emitAccountUpdates(conn, accountsGetResponse, &itemId, &token, &tenantId)

	assert.Nil(t, err)

	assert.Equal(t, 2, conn.sendMessageCount)

	var accountsChecked int
	for i, a := range conn.body {
		assert.Equal(t, token, conn.headers[i]["Authorization"])

		assert.Equal(t, "", conn.queue[i])

		assert.Equal(t, "account_update", conn.exchange[i])

		assert.Equal(t, "application/json", conn.contentType[i])

		account, ok := a.(Account)
		assert.True(t, ok)

		assert.Equal(t, "accountId", *account.AccountId)
		assert.Equal(t, "itemId", *account.ItemId)

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

	err := emitAccountUpdates(conn, accountsGetResponse, &itemId, &token, &tenantId)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to send one or more messages [failing for 1 test]\n", err.Error())

	assert.Equal(t, 2, conn.sendMessageCount)

	var accountsChecked int
	for _, a := range conn.body {

		account, ok := a.(Account)
		assert.True(t, ok)

		assert.Equal(t, "accountId", *account.AccountId)
		assert.Equal(t, "itemId", *account.ItemId)

		accountsChecked++
	}

	assert.Equal(t, 2, accountsChecked)
}

func createTestHandler(conn *TestConnector, middleware *TestMiddleware, plaidHandler *TestPlaidHandler) Handler {

	handler := NewHandler(conn, middleware, plaidHandler, nil)

	return handler
}

type TestConnector struct {
	body                 []interface{}
	headers              []map[string]interface{}
	contentType          []string
	queue                []string
	exchange             []string
	receiveMessagesCount int
	sendMessageCount     int
	failTimeCalled       int
}

func (conn *TestConnector) ReceiveMessages(queueName string, handler moneymakerrabbit.MessageHandlerFunc) {
	conn.receiveMessagesCount++
}

func (conn *TestConnector) SendMessage(body interface{}, headers map[string]interface{}, contentType string, queue string, exchange string) error {
	conn.body = append(conn.body, body)
	conn.headers = append(conn.headers, headers)
	conn.contentType = append(conn.contentType, contentType)
	conn.queue = append(conn.queue, queue)
	conn.exchange = append(conn.exchange, exchange)
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
