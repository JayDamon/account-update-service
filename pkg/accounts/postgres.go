package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type PostgresRepository struct {
	Conn *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		Conn: db,
	}
}

func (r *PostgresRepository) GetAccountsForUser(tenantId string) (*[]Account, error) {
	db := r.Conn

	statement := `SELECT a.account_id, a.friendly_name, a.name, a.mask, a.official_name, a.available_balance,
					a.current_balance, a.account_limit, a.is_primary_account, a.is_in_cash_flow, a.account_type, 
					a.account_sub_type, i.institution_id, i.institution_name
					FROM account a
					INNER JOIN item i ON i.item_id = a.item_id
					WHERE a.tenant_id = $1;`
	log.Printf("Running query: %s\n", statement)
	results, err := db.Query(statement, tenantId)
	if err != nil {
		log.Printf("error retrieving accounts for user %s: %s\n", tenantId, err)
		return nil, err
	}

	aa := make([]Account, 0)
	for results.Next() {
		var a Account
		err := results.Scan(&a.Id, &a.FriendlyName, &a.Name, &a.Mask, &a.OfficialName, &a.AvailableBalance,
			&a.CurrentBalance, &a.Limit, &a.IsPrimaryAccount, &a.IsInCashFlow, &a.AccountTypeName, &a.AccountSubTypeName,
			&a.InstitutionId, &a.InstitutionName)
		if err != nil {
			log.Printf("error caused while scanning results of query\n    Query: %s\n    Error: %s\n", statement, err)
		}
		aa = append(aa, a)
	}
	return &aa, nil
}
func (r *PostgresRepository) InsertNewAccounts(ai *AccountItem) (*AccountItem, error) {
	db := r.Conn
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		msg := "failed to begin database transaction"
		log.Println(msg)
		return nil, fmt.Errorf("%s: %s", msg, err)
	}
	defer tx.Rollback()

	err = saveItem(ai, ctx, tx)
	if err != nil {
		return nil, err
	}

	insertStmt := `INSERT INTO account (friendly_name, name, mask, plaid_id, item_id, official_name, 
                     available_balance, current_balance, starting_balance, account_limit, iso_currency_code, 
                     unofficial_currency_code, is_primary_account, is_in_cash_flow, account_type, 
                     account_sub_type, tenant_id) VALUES `
	var vals []interface{}

	vn := 0
	for _, a := range *ai.Accounts {
		insertStmt += "("
		// the high value here is the number of input parameters
		//  this will generate ($1, $2...) where the number is the current value of vn
		for i := 1; i <= 17; i++ {
			vn++
			insertStmt += fmt.Sprintf("$%d,", vn)
		}
		insertStmt = insertStmt[0:len(insertStmt)-1] + "),"

		available := newNullFloat64(a.AvailableBalance)
		current := newNullFloat64(a.CurrentBalance)
		limit := newNullFloat64(a.Limit)

		frNm := newNullString(a.FriendlyName)
		isoCurr := newNullString(a.OfficialCurrencyCode)
		unCurr := newNullString(a.UnofficialCurrencyCode)
		msk := newNullString(a.Mask)
		ofNm := newNullString(a.OfficialName)
		sbTp := newNullString(a.AccountSubTypeName)

		vals = append(vals,
			frNm,
			a.Name,
			msk,
			a.PlaidAccountId,
			ai.ItemId,
			ofNm,
			available,
			current,
			current,
			limit,
			isoCurr,
			unCurr,
			false,
			false,
			a.AccountTypeName,
			sbTp,
			ai.TenantId,
		)
	}

	insertStmt = insertStmt[0 : len(insertStmt)-1]
	log.Printf("Executing insert statmeent: %s\n values: %v\n", insertStmt, vals)
	_, err = tx.ExecContext(ctx, insertStmt, vals...)
	if err != nil {
		msg := "failed to insert new accounts"
		log.Printf("%s err: %s\n", msg, err)
		return nil, fmt.Errorf("%s: %s", msg, err)
	}
	log.Println("Successfully inserted new accounts")

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit changes: %s\n", err)
	}
	return nil, nil
}

func (r *PostgresRepository) UpdateTransactionName(a *Account) (*Account, error) {

	db := r.Conn
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		msg := "failed to begin database transaction"
		log.Println(msg)
		return nil, fmt.Errorf("%s: %s", msg, err)
	}
	defer tx.Rollback()

	updateStatement := `UPDATE account SET friendly_name = $1 WHERE account_id = $2;`
	vals := []interface{}{a.FriendlyName, a.Id}
	log.Printf("Running query: %s\n    vals: %s\n", updateStatement, vals)

	_, err = tx.ExecContext(ctx, updateStatement, vals...)
	if err != nil {
		msg := "error executing update statement for account name"
		log.Printf("%s err: %s\n", msg, err)
		return nil, err
	}

	return a, nil
}

func saveItem(ai *AccountItem, ctx context.Context, tx *sql.Tx) error {

	//selectStmt := `SELECT item_id FROM item WHERE item_id = $1;`
	//log.Printf("Checking if item already exists %s: %s\n", ai.ItemId, selectStmt)
	//row := tx.QueryRowContext(ctx, selectStmt, ai.ItemId)
	//// If row exists, return. No need to save
	//log.Printf("Select item error %s\n", row.Err())
	//if !errors.Is(row.Err(), sql.ErrNoRows) {
	//	log.Printf("Item already exists: %s\n", *ai.ItemId)
	//	return nil
	//}

	insertStmt := `INSERT INTO item (item_id, institution_id, institution_name, url, primary_color, logo, tenant_id)
					VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (item_id) DO NOTHING`

	iId := newNullString(ai.InstitutionId)
	in := newNullString(ai.InstitutionName)
	iUrl := newNullString(ai.Url)
	ipc := newNullString(ai.PrimaryColor)
	il := newNullString(ai.Logo)

	args := []interface{}{ai.ItemId, iId, in, iUrl, ipc, il, ai.TenantId}
	log.Printf("Running query: %s\n    with params: %s\n", insertStmt, args)
	_, err := tx.ExecContext(ctx, insertStmt, args...)
	if err != nil {
		msg := "failed to insert new accounts"
		log.Printf("%s err: %s\n", msg, err)
		return fmt.Errorf("%s: %s", msg, err)
	}
	log.Printf("Successfully created new item")
	return nil
}

func newNullFloat64(s *float32) sql.NullFloat64 {
	if s == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: float64(*s),
		Valid:   true,
	}
}

func newNullString(s *string) sql.NullString {
	if s == nil || len(*s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}
