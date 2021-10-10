package core

import (
	"context"
)

type Blog struct {
	BlogID  int
	FeedURL string
	SiteURL string
	Title   string
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
	Read(ctx context.Context, blogID int) (Blog, error)
	ReadAll(ctx context.Context) ([]Blog, error)
}
