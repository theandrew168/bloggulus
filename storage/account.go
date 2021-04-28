package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/model"
)

type Account interface {
	Create(ctx context.Context, account *model.Account) (*model.Account, error)
	Read(ctx context.Context, accountID int) (*model.Account, error)
	ReadByUsername(ctx context.Context, username string) (*model.Account, error)
	Delete(ctx context.Context, accountID int) error
}
