package model

import (
	"time"

	"uuid"
)

const (
	// Ensure that a post doesn't get synced more than once every SyncCooldown.
	SyncCooldown = 2 * time.Hour
)

type Blog struct {
	id           uuid.UUID
	feedURL      string
	siteURL      string
	title        string
	isPublic     bool
	syncedAt     time.Time
	etag         string
	lastModified string

	meta *Meta
}

func NewBlog(feedURL, siteURL, title string, syncedAt time.Time, etag, lastModified string) (*Blog, error) {
	blog := Blog{
		id:           uuid.New(),
		feedURL:      feedURL,
		siteURL:      siteURL,
		title:        title,
		isPublic:     false,
		syncedAt:     syncedAt,
		etag:         etag,
		lastModified: lastModified,

		meta: NewMeta(),
	}
	return &blog, nil
}

func LoadBlog(id uuid.UUID, feedURL, siteURL, title string, isPublic bool, syncedAt time.Time, etag, lastModified string, meta *Meta) *Blog {
	blog := Blog{
		id:           id,
		feedURL:      feedURL,
		siteURL:      siteURL,
		title:        title,
		isPublic:     isPublic,
		syncedAt:     syncedAt,
		etag:         etag,
		lastModified: lastModified,

		meta: meta,
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

func (b *Blog) IsPublic() bool {
	return b.isPublic
}

func (b *Blog) SetIsPublic(isPublic bool) error {
	b.isPublic = isPublic
	return nil
}

func (b *Blog) SyncedAt() time.Time {
	return b.syncedAt
}

func (b *Blog) SetSyncedAt(syncedAt time.Time) {
	b.syncedAt = syncedAt
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

func (b *Blog) Meta() *Meta {
	return b.meta
}

func (b *Blog) CanBeSynced(now time.Time) bool {
	return b.syncedAt.Add(SyncCooldown).Before(now)
}
