package core

import (
	"context"
	"time"
)

type Post struct {
	// fields known upfront
	URL     string    `json:"url"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
	Blog    Blog      `json:"blog"`

	// readonly (from database, after creation)
	ID   int      `json:"id"`
	Tags []string `json:"tags"`

	// used in sync process
	Body string `json:"-"`
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
