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
	ReadAll(ctx context.Context) ([]Blog, error)
}
