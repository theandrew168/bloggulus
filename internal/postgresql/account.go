package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type accountStorage struct {
	db *pgxpool.Pool
}

func NewAccountStorage(db *pgxpool.Pool) core.AccountStorage {
	s := accountStorage{
		db: db,
	}
	return &s
}

func (s *accountStorage) Create(ctx context.Context, account *core.Account) (*core.Account, error) {
	command := "INSERT INTO account (username, password, email) VALUES ($1, $2, $3) RETURNING account_id"
	row := s.db.QueryRow(ctx, command, account.Username, account.Password, account.Email)

	account.Verified = false
	err := row.Scan(&account.AccountID)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		// https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, core.ErrExist
			}
		}
		return nil, err
	}

	return account, nil
}

func (s *accountStorage) Read(ctx context.Context, accountID int) (*core.Account, error) {
	query := "SELECT * FROM account WHERE account_id = $1"
	row := s.db.QueryRow(ctx, query, accountID)

	var account core.Account
	err := row.Scan(&account.AccountID, &account.Username, &account.Password, &account.Email, &account.Verified)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *accountStorage) ReadByUsername(ctx context.Context, username string) (*core.Account, error) {
	query := "SELECT * FROM account WHERE username = $1"
	row := s.db.QueryRow(ctx, query, username)

	var account core.Account
	err := row.Scan(&account.AccountID, &account.Username, &account.Password, &account.Email, &account.Verified)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *accountStorage) Delete(ctx context.Context, accountID int) error {
	command := "DELETE FROM account WHERE account_id = $1"
	_, err := s.db.Exec(ctx, command, accountID)
	return err
}
