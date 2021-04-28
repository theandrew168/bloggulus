package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/model"
)

type Blog interface {
	Create(ctx context.Context, blog *model.Blog) (*model.Blog, error)
	Read(ctx context.Context, blogID int) (*model.Blog, error)
	ReadAll(ctx context.Context) ([]*model.Blog, error)
	ReadFollowedForUser(ctx context.Context, accountID int) ([]*model.Blog, error)
	ReadUnfollowedForUser(ctx context.Context, accountID int) ([]*model.Blog, error)
	Delete(ctx context.Context, blogID int) error
}
