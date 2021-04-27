package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
)

type Account interface {
	Create(ctx context.Context, account *models.Account) (*models.Account, error)
	Read(ctx context.Context, accountID int) (*models.Account, error)
	ReadByUsername(ctx context.Context, username string) (*models.Account, error)
	Delete(ctx context.Context, accountID int) error
}
