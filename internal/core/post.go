package core

import (
	"context"
	"time"
)

type Post struct {
	PostID  int
	URL     string
	Title   string
	Updated time.Time

	Blog Blog
}

type PostStorage interface {
	Create(ctx context.Context, post *Post) error
	Read(ctx context.Context, postID int) (Post, error)
	ReadAllByBlog(ctx context.Context, blogID int) ([]Post, error)
	ReadRecent(ctx context.Context, limit, offset int) ([]Post, error)
	ReadSearch(ctx context.Context, query string, limit, offset int) ([]Post, error)
	CountRecent(ctx context.Context) (int, error)
	CountSearch(ctx context.Context, query string) (int, error)
}
