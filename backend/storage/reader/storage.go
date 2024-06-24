package reader

import (
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type Storage struct {
	conn postgres.Conn

	article *ArticleStorage
}

func New(conn postgres.Conn) *Storage {
	s := Storage{
		conn: conn,

		article: NewArticleStorage(conn),
	}
	return &s
}

func (s *Storage) Article() *ArticleStorage {
	return s.article
}
