package query

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/value"
)

type Account struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	IsAdmin  bool      `db:"is_admin"`
}

type AccountQuery struct {
	conn postgres.Conn
}

func NewAccount(conn postgres.Conn) *AccountQuery {
	qry := AccountQuery{
		conn: conn,
	}
	return &qry
}

// Powers authentication middleware.
func (qry *AccountQuery) ReadBySessionTokenHash(sessionTokenHash value.TokenHash) (Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.is_admin
		FROM account
		INNER JOIN session
			ON session.account_id = account.id
		WHERE session.token_hash = $1;
	`

	rows, err := qry.conn.Query(context.Background(), stmt, sessionTokenHash.Value())
	if err != nil {
		return Account{}, err
	}

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Account])
	if err != nil {
		return Account{}, postgres.CheckReadError(err)
	}

	return account, nil
}

// Powers the accounts page (admin only).
func (qry *AccountQuery) List() ([]Account, error) {
	stmt := `
		SELECT
			account.id,
			account.username,
			account.is_admin
		FROM account;
	`

	rows, err := qry.conn.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Account])
	if err != nil {
		return nil, postgres.CheckListError(err)
	}

	return accounts, nil
}
