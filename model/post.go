package model

import (
	"time"
)

type Post struct {
	PostID  int
	BlogID  int
	URL     string
	Title   string
	Updated time.Time

	Blog Blog
}
