package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Post struct {
	id          uuid.UUID
	blogID      uuid.UUID
	url         string
	title       string
	content     string
	publishedAt time.Time

	createdAt time.Time
	updatedAt time.Time
}

func NewPost(blog *Blog, url, title, content string, publishedAt time.Time) (*Post, error) {
	now := timeutil.Now()
	post := Post{
		id:          uuid.New(),
		blogID:      blog.ID(),
		url:         url,
		title:       title,
		content:     content,
		publishedAt: publishedAt,

		createdAt: now,
		updatedAt: now,
	}
	return &post, nil
}

func LoadPost(id, blogID uuid.UUID, url, title, content string, publishedAt, createdAt, updatedAt time.Time) *Post {
	post := Post{
		id:          id,
		blogID:      blogID,
		url:         url,
		title:       title,
		content:     content,
		publishedAt: publishedAt,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &post
}

func (p *Post) ID() uuid.UUID {
	return p.id
}

func (p *Post) BlogID() uuid.UUID {
	return p.blogID
}

func (p *Post) URL() string {
	return p.url
}

func (p *Post) Title() string {
	return p.title
}

func (p *Post) SetTitle(title string) error {
	p.title = title
	return nil
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) SetContent(content string) error {
	p.content = content
	return nil
}

func (p *Post) PublishedAt() time.Time {
	return p.publishedAt
}

func (p *Post) SetPublishedAt(publishedAt time.Time) error {
	p.publishedAt = publishedAt
	return nil
}

func (p *Post) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Post) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Post) SetUpdatedAt(updatedAt time.Time) error {
	p.updatedAt = updatedAt
	return nil
}

func (p *Post) CheckDelete() error {
	return nil
}
