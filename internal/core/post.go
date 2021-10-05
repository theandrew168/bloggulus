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
	ReadRecent(ctx context.Context, n int) ([]Post, error)

	// select ...
	// from post
	// where content_index @@ websearch_to_tsquery('english',  $1);
	//	Search(ctx context.Context, query string) ([]Post, error)
}
