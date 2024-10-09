package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Page struct {
	id      uuid.UUID
	url     string
	title   string
	content string

	createdAt time.Time
	updatedAt time.Time
}

func NewPage(url, title, content string) (*Page, error) {
	now := timeutil.Now()
	page := Page{
		id:      uuid.New(),
		url:     url,
		title:   title,
		content: content,

		createdAt: now,
		updatedAt: now,
	}
	return &page, nil
}

func LoadPage(id uuid.UUID, url, title, content string, createdAt, updatedAt time.Time) *Page {
	page := Page{
		id:      id,
		url:     url,
		title:   title,
		content: content,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &page
}

func (p *Page) ID() uuid.UUID {
	return p.id
}

func (p *Page) URL() string {
	return p.url
}

func (p *Page) Title() string {
	return p.title
}

func (p *Page) Content() string {
	return p.content
}

func (p *Page) SetContent(content string) error {
	p.content = content
	return nil
}

func (p *Page) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Page) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Page) SetUpdatedAt(updatedAt time.Time) error {
	p.updatedAt = updatedAt
	return nil
}

func (p *Page) CheckDelete() error {
	return nil
}
