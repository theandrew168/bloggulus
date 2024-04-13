package reader

import (
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type Storage struct {
	conn postgres.Conn

	post *PostStorage
}

func New(conn postgres.Conn) *Storage {
	s := Storage{
		conn: conn,

		post: NewPostStorage(conn),
	}
	return &s
}

func (s *Storage) Post() *PostStorage {
	return s.post
}
