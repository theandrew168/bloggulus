package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
)

type Post interface {
	Create(ctx context.Context, post *models.Post) (*models.Post, error)
	Read(ctx context.Context, postID int) (*models.Post, error)
	ReadRecent(ctx context.Context, n int) ([]*models.Post, error)
	ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*models.Post, error)
	Delete(ctx context.Context, postID int) error
}
