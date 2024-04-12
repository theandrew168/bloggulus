package storage

import (
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/reader"
	"github.com/theandrew168/bloggulus/backend/domain/reader/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
)

// ensure PostStorage interface is satisfied
var _ storage.PostStorage = (*PostgresPostStorage)(nil)

type dbPost struct {
	Title       string    `db:"title"`
	URL         string    `db:"url"`
	BlogTitle   string    `db:"blog_title"`
	BlogURL     string    `db:"blog_url"`
	PublishedAt time.Time `db:"published_at"`
	Tags        []string  `db:"tags"`
}

func (p dbPost) unmarshal() (*reader.Post, error) {
	post := reader.LoadPost(
		p.Title,
		p.URL,
		p.BlogTitle,
		p.BlogURL,
		p.PublishedAt,
		p.Tags,
	)
	return post, nil
}

type PostgresPostStorage struct {
	conn postgres.Conn
}

func NewPostgresPostStorage(conn postgres.Conn) *PostgresPostStorage {
	s := PostgresPostStorage{
		conn: conn,
	}
	return &s
}

func (s *PostgresPostStorage) List(limit, offset int) ([]reader.Post, error) {
	return nil, nil
}

func (s *PostgresPostStorage) Search(query string, limit, offset int) ([]reader.Post, error) {
	return nil, nil
}
