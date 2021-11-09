package core

import (
	"context"
	"time"
)

type Post struct {
	// fields known upfront
	URL     string
	Title   string
	Updated time.Time
	Blog    Blog

	// readonly (from database, after creation)
	PostID int
	Tags   []string

	// used in sync process
	Body string
}

func NewPost(url, title string, updated time.Time, blog Blog) Post {
	post := Post{
		URL:     url,
		Title:   title,
		Updated: updated,
		Blog:    blog,
	}
	return post
}

type PostStorage interface {
	Create(ctx context.Context, post *Post) error
	ReadAllByBlog(ctx context.Context, blogID int) ([]Post, error)
	ReadRecent(ctx context.Context, limit, offset int) ([]Post, error)
	ReadSearch(ctx context.Context, query string, limit, offset int) ([]Post, error)
	CountRecent(ctx context.Context) (int, error)
	CountSearch(ctx context.Context, query string) (int, error)
}
