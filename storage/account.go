package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/model"
)

type Account struct {
	db *pgxpool.Pool
}

func NewAccount(db *pgxpool.Pool) *Account {
	s := Account{
		db: db,
	}
	return &s
}

func (s *Account) Create(ctx context.Context, account *model.Account) (*model.Account, error) {
	command := "INSERT INTO account (username, password, email) VALUES ($1, $2, $3) RETURNING account_id"
	row := s.db.QueryRow(ctx, command, account.Username, account.Password, account.Email)

	account.Verified = false
	err := row.Scan(&account.AccountID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Account) Read(ctx context.Context, accountID int) (*model.Account, error) {
	query := "SELECT * FROM account WHERE account_id = $1"
	row := s.db.QueryRow(ctx, query, accountID)

	var account model.Account
	err := row.Scan(&account.AccountID, &account.Username, &account.Password, &account.Email, &account.Verified)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Account) ReadByUsername(ctx context.Context, username string) (*model.Account, error) {
	query := "SELECT * FROM account WHERE username = $1"
	row := s.db.QueryRow(ctx, query, username)

	var account model.Account
	err := row.Scan(&account.AccountID, &account.Username, &account.Password, &account.Email, &account.Verified)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Account) Delete(ctx context.Context, accountID int) error {
	command := "DELETE FROM account WHERE account_id = $1"
	_, err := s.db.Exec(ctx, command, accountID)
	return err
}
