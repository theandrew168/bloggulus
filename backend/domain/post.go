package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID
	BlogID      uuid.UUID
	URL         string
	Title       string
	Content     string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPost(blog Blog, url, title, content string, publishedAt time.Time) Post {
	now := time.Now()
	post := Post{
		ID:          uuid.New(),
		BlogID:      blog.ID,
		URL:         url,
		Title:       title,
		Content:     content,
		PublishedAt: publishedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return post
}

func LoadPost(
	id uuid.UUID,
	blogID uuid.UUID,
	url string,
	title string,
	content string,
	publishedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Post {
	post := Post{
		ID:          id,
		BlogID:      blogID,
		URL:         url,
		Title:       title,
		Content:     content,
		PublishedAt: publishedAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	return post
}
