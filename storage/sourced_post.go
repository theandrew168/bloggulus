package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
)

type SourcedPost interface {
	ReadRecent(ctx context.Context, n int) ([]*models.SourcedPost, error)
	ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*models.SourcedPost, error)
}
