package model

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
    Create(ctx context.Context, blog *Blog) (*Blog, error)
    Read(ctx context.Context, blogID int) (*Blog, error)
    ReadByURL(ctx context.Context, feedURL string) (*Blog, error)
    ReadAll(ctx context.Context) ([]*Blog, error)
    ReadFollowedForUser(ctx context.Context, accountID int) ([]*Blog, error)
    ReadUnfollowedForUser(ctx context.Context, accountID int) ([]*Blog, error)
    Delete(ctx context.Context, blogID int) error
}
