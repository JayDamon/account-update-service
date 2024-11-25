package accounts

import (
	"github.com/plaid/plaid-go/plaid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertAccountResponseToAccountList(t *testing.T) {

	authGetResponse := createTestAccountsGetResponse()

	itemId := "itemId"
	tenantId := "tenantId"
	isNew := true

	accounts := convertAccountResponseToAccountList(authGetResponse, &itemId, &tenantId, &isNew)

	assert.NotEmpty(t, accounts, "Accounts array should not be nil or empty")
	assert.Len(t, *accounts, 2, "Accounts array should have 2 accounts, but was %n", len(*accounts))

	accountOneCounted := false
	accountTwoCounted := false

	for _, a := range *accounts {
		assert.NotNil(t, a.ItemId)
		assert.Equal(t, *a.ItemId, "itemId")
		assert.Equal(t, *a.TenantId, "tenantId")
		assert.True(t, *a.IsNew)
		if *a.PlaidAccountId == "accountOne" {
			assert.Equal(t, *a.Name, "testName")
			accountOneCounted = true
		} else {
			assert.Equal(t, *a.PlaidAccountId, "accountTwo")
			assert.Equal(t, *a.Name, "testNameTwo")
			accountTwoCounted = true
		}
	}

	assert.True(t, accountOneCounted, "AccountOneCounted should be true")
	assert.True(t, accountTwoCounted, "AccountTwoCounted should be true")
}

func TestConvertAccountBaseToAccount_HappyPath(t *testing.T) {

	testName := "testName"
	officialName := "officialName"
	testId := "accountOne"
	available := 22.3
	current := 31.4
	limit := 500
	accType := plaid.ACCOUNTTYPE_INVESTMENT
	accSubType := plaid.ACCOUNTSUBTYPE__401K

	accountBase := createSimpleAccountBase(testName, officialName, available, current, limit, accType, accSubType)

	account := convertAccountBaseToAccount(*accountBase)

	assert.NotEmpty(t, account, "Account must not be nil or empty")
	assert.Equal(t, testId, *account.PlaidAccountId)
	assert.Equal(t, testName, *account.Name)
	assert.Nil(t, account.ItemId)
	assert.Equal(t, officialName, *account.OfficialName)
	assert.Equal(t, float32(available), *account.AvailableBalance)
	assert.Equal(t, float32(current), *account.CurrentBalance)
	assert.Equal(t, float32(limit), *account.Limit)
	assert.Equal(t, string(accType), *account.AccountTypeName)
	assert.Equal(t, string(accSubType), *account.AccountSubTypeName)
}

func TestConvertAccountBaseToAccount_NilOfficialName(t *testing.T) {

	testName := "testName"
	testId := "id"
	mask := "0000"
	currencyCode := "USD"
	var available float32 = 22.3
	var current float32 = 31.4
	var limit float32 = 500
	accType := plaid.ACCOUNTTYPE_INVESTMENT
	accSubType := plaid.ACCOUNTSUBTYPE__401K

	accountBase := createNullableAccountBase(
		testName,
		*plaid.NewNullableString(nil),
		*plaid.NewNullableString(&mask),
		*plaid.NewNullableFloat32(&available),
		*plaid.NewNullableFloat32(&current),
		*plaid.NewNullableFloat32(&limit),
		*plaid.NewNullableString(&currencyCode),
		accType,
		*plaid.NewNullableAccountSubtype(&accSubType),
	)

	account := convertAccountBaseToAccount(*accountBase)

	assert.NotEmpty(t, account, "Account must not be nil or empty")
	assert.Equal(t, testId, *account.PlaidAccountId)
	assert.Equal(t, testName, *account.Name)
	assert.Nil(t, account.ItemId)
	assert.Nil(t, account.OfficialName)
	assert.Equal(t, available, *account.AvailableBalance)
	assert.Equal(t, current, *account.CurrentBalance)
	assert.Equal(t, limit, *account.Limit)
	assert.Equal(t, string(accType), *account.AccountTypeName)
	assert.Equal(t, string(accSubType), *account.AccountSubTypeName)
}

func TestConvertAccountBaseToAccount_NilAvailableBalance(t *testing.T) {

	testName := "testName"
	testId := "id"
	mask := "0000"
	currencyCode := "USD"
	var current float32 = 31.4
	var limit float32 = 500
	accType := plaid.ACCOUNTTYPE_INVESTMENT
	accSubType := plaid.ACCOUNTSUBTYPE__401K

	accountBase := createNullableAccountBase(
		testName,
		*plaid.NewNullableString(&testName),
		*plaid.NewNullableString(&mask),
		*plaid.NewNullableFloat32(nil),
		*plaid.NewNullableFloat32(&current),
		*plaid.NewNullableFloat32(&limit),
		*plaid.NewNullableString(&currencyCode),
		accType,
		*plaid.NewNullableAccountSubtype(&accSubType),
	)

	account := convertAccountBaseToAccount(*accountBase)

	assert.NotEmpty(t, account, "Account must not be nil or empty")
	assert.Equal(t, testId, *account.PlaidAccountId)
	assert.Equal(t, testName, *account.Name)
	assert.Nil(t, account.ItemId)
	assert.Equal(t, testName, *account.OfficialName)
	assert.Nil(t, account.AvailableBalance)
	assert.Equal(t, current, *account.CurrentBalance)
	assert.Equal(t, limit, *account.Limit)
	assert.Equal(t, string(accType), *account.AccountTypeName)
	assert.Equal(t, string(accSubType), *account.AccountSubTypeName)
}

func TestConvertAccountBaseToAccount_NilCurrentBalance(t *testing.T) {

	testName := "testName"
	testId := "id"
	mask := "0000"
	currencyCode := "USD"
	var available float32 = 22.3
	var limit float32 = 500
	accType := plaid.ACCOUNTTYPE_INVESTMENT
	accSubType := plaid.ACCOUNTSUBTYPE__401K

	accountBase := createNullableAccountBase(
		testName,
		*plaid.NewNullableString(&testName),
		*plaid.NewNullableString(&mask),
		*plaid.NewNullableFloat32(&available),
		*plaid.NewNullableFloat32(nil),
		*plaid.NewNullableFloat32(&limit),
		*plaid.NewNullableString(&currencyCode),
		accType,
		*plaid.NewNullableAccountSubtype(&accSubType),
	)

	account := convertAccountBaseToAccount(*accountBase)

	assert.NotEmpty(t, account, "Account must not be nil or empty")
	assert.Equal(t, testId, *account.PlaidAccountId)
	assert.Equal(t, testName, *account.Name)
	assert.Nil(t, account.ItemId)
	assert.Equal(t, testName, *account.OfficialName)
	assert.Equal(t, available, *account.AvailableBalance)
	assert.Nil(t, account.CurrentBalance)
	assert.Equal(t, limit, *account.Limit)
	assert.Equal(t, string(accType), *account.AccountTypeName)
	assert.Equal(t, string(accSubType), *account.AccountSubTypeName)
}

func TestConvertAccountBaseToAccount_NilLimit(t *testing.T) {

	testName := "testName"
	testId := "id"
	mask := "0000"
	currencyCode := "USD"
	var available float32 = 22.3
	var current float32 = 31.4
	accType := plaid.ACCOUNTTYPE_INVESTMENT
	accSubType := plaid.ACCOUNTSUBTYPE__401K

	accountBase := createNullableAccountBase(
		testName,
		*plaid.NewNullableString(&testName),
		*plaid.NewNullableString(&mask),
		*plaid.NewNullableFloat32(&available),
		*plaid.NewNullableFloat32(&current),
		*plaid.NewNullableFloat32(nil),
		*plaid.NewNullableString(&currencyCode),
		accType,
		*plaid.NewNullableAccountSubtype(&accSubType),
	)

	account := convertAccountBaseToAccount(*accountBase)

	assert.NotEmpty(t, account, "Account must not be nil or empty")
	assert.Equal(t, testId, *account.PlaidAccountId)
	assert.Equal(t, testName, *account.Name)
	assert.Nil(t, account.ItemId)
	assert.Equal(t, testName, *account.OfficialName)
	assert.Equal(t, available, *account.AvailableBalance)
	assert.Equal(t, current, *account.CurrentBalance)
	assert.Nil(t, account.Limit)
	assert.Equal(t, string(accType), *account.AccountTypeName)
	assert.Equal(t, string(accSubType), *account.AccountSubTypeName)
}

func TestConvertAccountBaseToAccount_NilAccountSubType(t *testing.T) {

	testName := "testName"
	testId := "id"
	mask := "0000"
	currencyCode := "USD"
	var available float32 = 22.3
	var current float32 = 31.4
	var limit float32 = 500
	accType := plaid.ACCOUNTTYPE_INVESTMENT

	accountBase := createNullableAccountBase(
		testName,
		*plaid.NewNullableString(&testName),
		*plaid.NewNullableString(&mask),
		*plaid.NewNullableFloat32(&available),
		*plaid.NewNullableFloat32(&current),
		*plaid.NewNullableFloat32(&limit),
		*plaid.NewNullableString(&currencyCode),
		accType,
		*plaid.NewNullableAccountSubtype(nil),
	)

	account := convertAccountBaseToAccount(*accountBase)

	assert.NotEmpty(t, account, "Account must not be nil or empty")
	assert.Equal(t, testId, *account.PlaidAccountId)
	assert.Equal(t, testName, *account.Name)
	assert.Nil(t, account.ItemId)
	assert.Equal(t, testName, *account.OfficialName)
	assert.Equal(t, available, *account.AvailableBalance)
	assert.Equal(t, current, *account.CurrentBalance)
	assert.Equal(t, limit, *account.Limit)
	assert.Equal(t, string(accType), *account.AccountTypeName)
	assert.Nil(t, account.AccountSubTypeName)
}

func createTestAccountsGetResponse() *plaid.AccountsGetResponse {

	accountBalance := createAccountBalance(22.3, 33.1, 100.0, "USD")
	accountBalanceTwo := createAccountBalance(22.3, 33.1, 100.0, "USD")

	accountOne401A := createAccountBase("accountOne", "testName", accountBalance, "0000", "testName", plaid.ACCOUNTTYPE_INVESTMENT, plaid.ACCOUNTSUBTYPE__401A)
	account401K := createAccountBase("accountTwo", "testNameTwo", accountBalanceTwo, "0000", "testNameTwo", plaid.ACCOUNTTYPE_INVESTMENT, plaid.ACCOUNTSUBTYPE__401K)

	accounts := []plaid.AccountBase{
		*accountOne401A,
		*account401K,
	}

	item := plaid.NewItemWithDefaults()
	item.ItemId = "itemId"

	return plaid.NewAccountsGetResponse(
		accounts,
		*item,
		"test")
}

func createSimpleAccountBase(
	name string,
	officialName string,
	available float64,
	current float64,
	limit int,
	accType plaid.AccountType,
	subType plaid.AccountSubtype) *plaid.AccountBase {
	balance := createAccountBalance(float32(available), float32(current), float32(limit), "USD")
	accountBase := createAccountBase("accountOne", name, balance, "0000", officialName, accType, subType)
	return accountBase
}

func createNullableAccountBase(
	name string,
	officialName plaid.NullableString,
	mask plaid.NullableString,
	available plaid.NullableFloat32,
	current plaid.NullableFloat32,
	limit plaid.NullableFloat32,
	currencyCode plaid.NullableString,
	accType plaid.AccountType,
	subType plaid.NullableAccountSubtype) *plaid.AccountBase {

	balance := plaid.NewAccountBalance(
		available,
		current,
		limit,
		currencyCode,
		currencyCode)

	accountBase := plaid.NewAccountBase(
		"id",
		*balance,
		mask,
		name,
		officialName,
		accType,
		subType,
	)
	return accountBase
}

func createAccountBalance(available float32, current float32, limit float32, currencyCode string) *plaid.AccountBalance {
	accountBalance := plaid.NewAccountBalance(
		*plaid.NewNullableFloat32(&available),
		*plaid.NewNullableFloat32(&current),
		*plaid.NewNullableFloat32(&limit),
		*plaid.NewNullableString(&currencyCode),
		*plaid.NewNullableString(nil))
	return accountBalance
}

func createAccountBase(accountId string, accountName string, accountBalance *plaid.AccountBalance, mask string, officialName string, accType plaid.AccountType, subType plaid.AccountSubtype) *plaid.AccountBase {

	accountOne := plaid.NewAccountBase(
		accountId,
		*accountBalance,
		*plaid.NewNullableString(&mask),
		accountName,
		*plaid.NewNullableString(&officialName),
		accType,
		*plaid.NewNullableAccountSubtype(&subType),
	)
	return accountOne
}
