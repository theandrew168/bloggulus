package postgresql_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestCreateBlog(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.CreateBlog(storage, t)
}

func TestCreateBlogAlreadyExists(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.CreateBlogAlreadyExists(storage, t)
}

func TestReadBlog(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.ReadBlog(storage, t)
}

func TestReadBlogs(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	storage := postgresql.NewStorage(conn)
	test.ReadBlogs(storage, t)
}
