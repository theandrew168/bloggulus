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
	test.BlogCreate(t, storage)
}

func TestBlogCreateExists(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.BlogCreateExists(t, storage)
}

func TestBlogReadAll(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.BlogReadAll(t, storage)
}
