package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
)

type Blog interface {
	Create(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	Read(ctx context.Context, blogID int) (*models.Blog, error)
	ReadAll(ctx context.Context) ([]*models.Blog, error)
	ReadFollowedForUser(ctx context.Context, accountID int) ([]*models.Blog, error)
	ReadUnfollowedForUser(ctx context.Context, accountID int) ([]*models.Blog, error)
	Delete(ctx context.Context, blogID int) error
}
