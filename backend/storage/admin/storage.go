package admin

import (
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type Storage struct {
	conn postgres.Conn

	blog *BlogStorage
	post *PostStorage
	tag  *TagStorage
}

func New(conn postgres.Conn) *Storage {
	s := Storage{
		conn: conn,

		blog: NewBlogStorage(conn),
		post: NewPostStorage(conn),
		tag:  NewTagStorage(conn),
	}
	return &s
}

func (s *Storage) Blog() *BlogStorage {
	return s.blog
}

func (s *Storage) Post() *PostStorage {
	return s.post
}

func (s *Storage) Tag() *TagStorage {
	return s.tag
}
