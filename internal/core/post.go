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
	ReadAllByBlog(ctx context.Context, blogID int) ([]Post, error)
	ReadRecent(ctx context.Context, n int) ([]Post, error)
	ReadRecentByAccount(ctx context.Context, accountID int, n int) ([]Post, error)
	Delete(ctx context.Context, postID int) error

	// select ...
	// from post
	// where body_index @@ websearch_to_tsquery('english',  $1);
//	Search(ctx context.Context, query string) ([]Post, error)
}
