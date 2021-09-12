package postgresql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type accountStorage struct {
	conn *pgxpool.Pool
}

func NewAccountStorage(conn *pgxpool.Pool) core.AccountStorage {
	s := accountStorage{
		conn: conn,
	}
	return &s
}

func (s *accountStorage) Create(ctx context.Context, account *core.Account) error {
	stmt := `
		INSERT INTO account
			(username, password, email, verified)
		VALUES
			($1, $2, $3, $4)
		RETURNING account_id`
	err := pgxscan.Get(ctx, s.conn, account, stmt,
		account.Username,
		account.Password,
		account.Email,
		account.Verified)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return core.ErrExist
			}
		}
		return err
	}

	return nil
}

func (s *accountStorage) Read(ctx context.Context, accountID int) (core.Account, error) {
	stmt := "SELECT * FROM account WHERE account_id = $1"
	row := s.conn.QueryRow(ctx, stmt, accountID)

	var account core.Account
	err := row.Scan(
		&account.AccountID,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.Verified)
	if err != nil {
		return core.Account{}, err
	}

	return account, nil
}

func (s *accountStorage) ReadByUsername(ctx context.Context, username string) (core.Account, error) {
	stmt := "SELECT * FROM account WHERE username = $1"
	row := s.conn.QueryRow(ctx, stmt, username)

	var account core.Account
	err := row.Scan(
		&account.AccountID,
		&account.Username,
		&account.Password,
		&account.Email,
		&account.Verified)
	if err != nil {
		return core.Account{}, err
	}

	return account, nil
}

func (s *accountStorage) Delete(ctx context.Context, accountID int) error {
	stmt := "DELETE FROM account WHERE account_id = $1"
	_, err := s.conn.Exec(ctx, stmt, accountID)
	return err
}
