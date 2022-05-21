package bloggulus

import (
	"time"
)

type Post struct {
	// fields known upfront
	URL     string    `json:"url"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
	Blog    Blog      `json:"blog"`

	// readonly (from database, after creation)
	ID   int      `json:"id"`
	Tags []string `json:"tags"`

	// used in sync process
	Body string `json:"-"`
}

// NewPost creates a new Post struct.
func NewPost(url, title string, updated time.Time, blog Blog) Post {
	post := Post{
		URL:     url,
		Title:   title,
		Updated: updated,
		Blog:    blog,
	}
	return post
}
