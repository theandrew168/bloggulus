package core

import (
	"context"
)

type Blog struct {
	BlogID  int
	FeedURL string
	SiteURL string
	Title   string
}

type BlogStorage interface {
	Create(ctx context.Context, blog *Blog) error
	Read(ctx context.Context, blogID int) (Blog, error)
	ReadByURL(ctx context.Context, feedURL string) (Blog, error)
	ReadAll(ctx context.Context) ([]Blog, error)
	ReadFollowedByAccount(ctx context.Context, accountID int) ([]Blog, error)
	ReadUnfollowedByAccount(ctx context.Context, accountID int) ([]Blog, error)
	Delete(ctx context.Context, blogID int) error
}
