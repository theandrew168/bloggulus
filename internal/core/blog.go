package core

import (
	"context"
)

type Blog struct {
	FeedURL string
	SiteURL string
	Title   string

	// readonly (from database, after creation)
	BlogID  int
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
	Create(ctx context.Context, blog *Blog) error
	ReadAll(ctx context.Context) ([]Blog, error)
}
