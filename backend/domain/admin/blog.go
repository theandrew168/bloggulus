package admin

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	id           uuid.UUID
	feedURL      string
	siteURL      string
	title        string
	etag         string
	lastModified string
	syncedAt     time.Time

	createdAt time.Time
	updatedAt time.Time
}

func NewBlog(feedURL, siteURL, title, etag, lastModified string, syncedAt time.Time) (*Blog, error) {
	now := time.Now().UTC()
	blog := Blog{
		id:           uuid.New(),
		feedURL:      feedURL,
		siteURL:      siteURL,
		title:        title,
		etag:         etag,
		lastModified: lastModified,
		syncedAt:     syncedAt,

		createdAt: now,
		updatedAt: now,
	}
	return &blog, nil
}

func LoadBlog(id uuid.UUID, feedURL, siteURL, title, etag, lastModified string, syncedAt, createdAt, updatedAt time.Time) *Blog {
	blog := Blog{
		id:           id,
		feedURL:      feedURL,
		siteURL:      siteURL,
		title:        title,
		etag:         etag,
		lastModified: lastModified,
		syncedAt:     syncedAt,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &blog
}

func (b *Blog) ID() uuid.UUID {
	return b.id
}

func (b *Blog) FeedURL() string {
	return b.feedURL
}

func (b *Blog) SiteURL() string {
	return b.siteURL
}

func (b *Blog) Title() string {
	return b.title
}

func (b *Blog) ETag() string {
	return b.etag
}

func (b *Blog) SetETag(etag string) error {
	b.etag = etag
	return nil
}

func (b *Blog) LastModified() string {
	return b.lastModified
}

func (b *Blog) SetLastModified(lastModified string) error {
	b.lastModified = lastModified
	return nil
}

func (b *Blog) SyncedAt() time.Time {
	return b.syncedAt
}

func (b *Blog) SetSyncedAt(syncedAt time.Time) {
	b.syncedAt = syncedAt
}

func (b *Blog) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Blog) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *Blog) SetUpdatedAt(updatedAt time.Time) error {
	b.updatedAt = updatedAt
	return nil
}

func (b *Blog) CheckDelete() error {
	return nil
}
