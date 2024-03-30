package admin

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
	SyncedAt     time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBlog(feedURL, siteURL, title, etag, lastModified string, syncedAt time.Time) Blog {
	now := time.Now()
	blog := Blog{
		ID:           uuid.New(),
		FeedURL:      feedURL,
		SiteURL:      siteURL,
		Title:        title,
		ETag:         etag,
		LastModified: lastModified,
		SyncedAt:     syncedAt,

		CreatedAt: now,
		UpdatedAt: now,
	}
	return blog
}
