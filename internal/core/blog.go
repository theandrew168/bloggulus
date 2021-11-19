package core

import (
	"context"
)

type Blog struct {
	FeedURL string
	SiteURL string
	Title   string

	// readonly (from database, after creation)
	ID int
}

func NewBlog(feedURL, siteURL, title string) Blog {
	blog := Blog{
		FeedURL: feedURL,
		SiteURL: siteURL,
		Title:   title,
	}
	return blog
}

type BlogStorage interface {
	BlogCreate(ctx context.Context, blog *Blog) error
	BlogReadAll(ctx context.Context) ([]Blog, error)
}
