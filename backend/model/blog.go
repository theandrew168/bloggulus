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

type NewBlogParams struct {
	FeedURL      string
	SiteURL      string
	Title        string
	SyncedAt     time.Time
	ETag         string
	LastModified string
}

func NewBlog(params NewBlogParams) (*Blog, error) {
	blog := Blog{
		id:           uuid.New(),
		feedURL:      params.FeedURL,
		siteURL:      params.SiteURL,
		title:        params.Title,
		isPublic:     false,
		syncedAt:     params.SyncedAt,
		etag:         params.ETag,
		lastModified: params.LastModified,
		meta:         NewMeta(),
	}
	return &blog, nil
}

type LoadBlogParams struct {
	ID           uuid.UUID
	FeedURL      string
	SiteURL      string
	Title        string
	IsPublic     bool
	SyncedAt     time.Time
	ETag         string
	LastModified string
	Meta         *Meta
}

func LoadBlog(params LoadBlogParams) *Blog {
	blog := Blog{
		id:           params.ID,
		feedURL:      params.FeedURL,
		siteURL:      params.SiteURL,
		title:        params.Title,
		isPublic:     params.IsPublic,
		syncedAt:     params.SyncedAt,
		etag:         params.ETag,
		lastModified: params.LastModified,
		meta:         params.Meta,
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
