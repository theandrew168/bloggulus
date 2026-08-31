package model

import (
	"time"
	"uuid"
)

type Post struct {
	id          uuid.UUID
	blogID      uuid.UUID
	url         string
	title       string
	publishedAt time.Time
	content     string

	meta *Meta
}

func NewPost(blog *Blog, url, title string, publishedAt time.Time, content string) (*Post, error) {
	post := Post{
		id:          uuid.New(),
		blogID:      blog.ID(),
		url:         url,
		title:       title,
		publishedAt: publishedAt,
		content:     content,

		meta: NewMeta(),
	}
	return &post, nil
}

func LoadPost(id, blogID uuid.UUID, url, title string, publishedAt time.Time, content string, meta *Meta) *Post {
	post := Post{
		id:          id,
		blogID:      blogID,
		url:         url,
		title:       title,
		publishedAt: publishedAt,
		content:     content,

		meta: meta,
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

func (p *Post) PublishedAt() time.Time {
	return p.publishedAt
}

func (p *Post) SetPublishedAt(publishedAt time.Time) error {
	p.publishedAt = publishedAt
	return nil
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) SetContent(content string) error {
	p.content = content
	return nil
}

func (p *Post) Meta() *Meta {
	return p.meta
}
