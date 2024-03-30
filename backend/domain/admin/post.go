package admin

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID
	BlogID      uuid.UUID
	URL         string
	Title       string
	Contents    string
	PublishedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPost(blog Blog, url, title, content string, publishedAt time.Time) Post {
	now := time.Now()
	post := Post{
		ID:          uuid.New(),
		BlogID:      blog.ID,
		URL:         url,
		Title:       title,
		Contents:    content,
		PublishedAt: publishedAt,

		CreatedAt: now,
		UpdatedAt: now,
	}
	return post
}
