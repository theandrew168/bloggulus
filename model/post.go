package model

import (
	"context"
	"time"
)

type Post struct {
	PostID  int
	BlogID  int
	URL     string
	Title   string
	Updated time.Time

	Blog Blog
}

type PostStorage interface {
    Create(ctx context.Context, post *Post) (*Post, error)
    Read(ctx context.Context, postID int) (*Post, error)
    ReadRecent(ctx context.Context, n int) ([]*Post, error)
    ReadRecentForUser(ctx context.Context, accountID int, n int) ([]*Post, error)
    Delete(ctx context.Context, postID int) error
}
