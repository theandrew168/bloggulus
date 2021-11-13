package postgresql_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestBlogCreate(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.BlogCreate(storage, t)
}

func TestBlogCreateExists(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.BlogCreateExists(storage, t)
}

func TestBlogReadAll(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.BlogReadAll(storage, t)
}
