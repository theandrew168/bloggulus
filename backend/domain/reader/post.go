package reader

import (
	"time"
)

type Post struct {
	Title     string
	URL       string
	BlogTitle string
	BlogURL   string
	Updated   time.Time
	Tags      []string
}
