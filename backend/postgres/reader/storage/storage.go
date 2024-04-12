package storage

import (
	"github.com/theandrew168/bloggulus/backend/domain/reader/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// ensure Storage interface is satisfied
var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	conn postgres.Conn

	post *PostgresPostStorage
}

func New(conn postgres.Conn) *Storage {
	s := Storage{
		conn: conn,

		post: NewPostgresPostStorage(conn),
	}
	return &s
}

func (s *Storage) Post() storage.PostStorage {
	return s.post
}
