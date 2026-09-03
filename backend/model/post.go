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

type NewPostParams struct {
	Blog        *Blog
	URL         string
	Title       string
	PublishedAt time.Time
	Content     string
}

func NewPost(params NewPostParams) (*Post, error) {
	post := Post{
		id:          uuid.New(),
		blogID:      params.Blog.ID(),
		url:         params.URL,
		title:       params.Title,
		publishedAt: params.PublishedAt,
		content:     params.Content,
		meta:        NewMeta(),
	}
	return &post, nil
}

type LoadPostParams struct {
	ID          uuid.UUID
	BlogID      uuid.UUID
	URL         string
	Title       string
	PublishedAt time.Time
	Content     string
	Meta        *Meta
}

func LoadPost(params LoadPostParams) *Post {
	post := Post{
		id:          params.ID,
		blogID:      params.BlogID,
		url:         params.URL,
		title:       params.Title,
		publishedAt: params.PublishedAt,
		content:     params.Content,
		meta:        params.Meta,
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
