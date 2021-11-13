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
	test.PostCreate(storage, t)
}

func TestPostCreateExists(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostCreateExists(storage, t)
}

func TestPostReadAllByBlog(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostReadAllByBlog(storage, t)
}

func TestPostReadRecent(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostReadRecent(storage, t)
}

func TestPostReadSearch(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostReadSearch(storage, t)
}

func TestPostCountRecent(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostCountRecent(storage, t)
}

func TestPostCountSearch(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.PostCountSearch(storage, t)
}
