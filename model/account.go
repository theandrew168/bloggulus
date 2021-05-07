package model

import (
	"context"
)

type Account struct {
	AccountID int
	Username  string
	Password  string
	Email     string
	Verified  bool
}

type AccountStorage interface {
    Create(ctx context.Context, account *Account) (*Account, error)
    Read(ctx context.Context, accountID int) (*Account, error)
    ReadByUsername(ctx context.Context, username string) (*Account, error)
    Delete(ctx context.Context, accountID int) error
}
