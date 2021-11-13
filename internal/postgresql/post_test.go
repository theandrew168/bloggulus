package postgresql_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestPostCreate(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostCreate(t, storage)
}

func TestPostCreateExists(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostCreateExists(t, storage)
}

func TestPostReadAllByBlog(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostReadAllByBlog(t, storage)
}

func TestPostReadRecent(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostReadRecent(t, storage)
}

func TestPostReadSearch(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostReadSearch(t, storage)
}

func TestPostCountRecent(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostCountRecent(t, storage)
}

func TestPostCountSearch(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostCountSearch(t, storage)
}
