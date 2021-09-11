package core

import (
	"context"
	"time"
)

type Post struct {
	PostID  int
	URL     string
	Title   string
	Author  string
	Body    string
	Updated time.Time

	Blog Blog
}

type PostStorage interface {
	Create(ctx context.Context, post *Post) error
	Read(ctx context.Context, postID int) (Post, error)
	ReadRecent(ctx context.Context, n int) ([]Post, error)
	ReadRecentForUser(ctx context.Context, accountID int, n int) ([]Post, error)
	Delete(ctx context.Context, postID int) error
}
