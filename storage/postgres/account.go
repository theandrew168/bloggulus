package postgres

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type accountStorage struct {
	db *pgxpool.Pool
}

func NewAccountStorage(db *pgxpool.Pool) storage.Account {
	return &accountStorage{
		db: db,
	}
}

func (s *accountStorage) Create(ctx context.Context, account *models.Account) (*models.Account, error) {
	command := "INSERT INTO account (username, password, email) VALUES ($1, $2, $3) RETURNING account_id"
	row := s.db.QueryRow(ctx, command, account.Username, account.Password, account.Email)

	account.Verified = false
	err := row.Scan(&account.AccountID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountStorage) Read(ctx context.Context, accountID int) (*models.Account, error) {
	query := "SELECT * FROM account WHERE account_id = $1"
	row := s.db.QueryRow(ctx, query, accountID)

	var account models.Account
	err := row.Scan(&account.AccountID, &account.Username, &account.Password, &account.Email, &account.Verified)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *accountStorage) ReadByUsername(ctx context.Context, username string) (*models.Account, error) {
	query := "SELECT * FROM account WHERE username = $1"
	row := s.db.QueryRow(ctx, query, username)

	var account models.Account
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
