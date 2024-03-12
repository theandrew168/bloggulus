package domain

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID           uuid.UUID
	FeedURL      string
	SiteURL      string
	Title        string
	ETag         string
	LastModified string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewBlog(feedURL, siteURL, title, etag, lastModified string) Blog {
	now := time.Now()
	blog := Blog{
		ID:           uuid.New(),
		FeedURL:      feedURL,
		SiteURL:      siteURL,
		Title:        title,
		ETag:         etag,
		LastModified: lastModified,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return blog
}

func LoadBlog(
	id uuid.UUID,
	feedURL string,
	siteURL string,
	title string,
	etag string,
	lastModified string,
	createdAt time.Time,
	updatedAt time.Time,
) Blog {
	blog := Blog{
		ID:           id,
		FeedURL:      feedURL,
		SiteURL:      siteURL,
		Title:        title,
		ETag:         etag,
		LastModified: lastModified,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
	return blog
}
