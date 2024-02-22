package accounts

import "testing"

func TestEmitAccountUpdates_SuccessPath(t *testing.T) {

	handler := createTestHandler()

	accounts := createTestAccountsGetResponse()
	token := "token"

	handler.emitAccountUpdates(accounts, &token)


}

func createTestHandler() *Handler {

	conn := &TestConnection{}

	handler := NewHandler(conn, nil, nil)

	return handler
}

type TestConnection struct {

}

func (conn *TestConnection) ReceiveMessages(queueName string, handler MessageHandlerFunc) {

}

func (conn *TestConnection) SendMessage(body interface{}, headers map[string]interface{}, contentType string, queue string) {

}