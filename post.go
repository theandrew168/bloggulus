package bloggulus

import (
	"time"
)

type Post struct {
	// fields known upfront
	URL     string    `json:"url"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`

	// used in sync process
	Body string `json:"-"`

	// belongs to a single blog
	Blog Blog `json:"blog"`

	// readonly (from database, after creation)
	ID   int      `json:"id"`
	Tags []string `json:"tags"`
}

// NewPost creates a new Post struct.
func NewPost(url, title string, updated time.Time, body string, blog Blog) Post {
	post := Post{
		URL:     url,
		Title:   title,
		Updated: updated,
		Body:    body,
		Blog:    blog,
	}
	return post
}
