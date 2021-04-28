package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/model"
)

type Post interface {
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	Read(ctx context.Context, postID int) (*model.Post, error)
	ReadRecent(ctx context.Context, n int) ([]*model.Post, error)
	ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*model.Post, error)
	Delete(ctx context.Context, postID int) error
}
