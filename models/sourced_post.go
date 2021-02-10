package models

import (
	"time"
)

type SourcedPost struct {
	URL       string
	Title     string
	Updated   time.Time
	BlogTitle string
}
