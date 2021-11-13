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
	PostCreate(ctx context.Context, post *Post) error
	PostReadAllByBlog(ctx context.Context, blogID int) ([]Post, error)
	PostReadRecent(ctx context.Context, limit, offset int) ([]Post, error)
	PostReadSearch(ctx context.Context, query string, limit, offset int) ([]Post, error)
	PostCountRecent(ctx context.Context) (int, error)
	PostCountSearch(ctx context.Context, query string) (int, error)
}
