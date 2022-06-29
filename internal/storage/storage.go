package storage

import (
	"time"

	"github.com/theandrew168/bloggulus/internal/database"
)

// default query timeout
const timeout = 10 * time.Second

type Storage struct {
	db database.Conn

	Blog *Blog
	Post *Post
}

func New(db database.Conn) *Storage {
	s := Storage{
		db: db,

		Blog: NewBlog(db),
		Post: NewPost(db),
	}
	return &s
}
