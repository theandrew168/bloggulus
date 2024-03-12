package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewMockBlog()
		err := store.Blog.Create(blog)
		if err != nil {
			t.Fatal(err)
		}

		return test.ErrSkipCommit
	})
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateMockBlog(t, store)

		// attempt to create the same blog again
		err := store.Blog.Create(blog)
		if !errors.Is(err, storage.ErrConflict) {
			t.Fatal("duplicate blog should return an error")
		}

		return test.ErrSkipCommit
	})
}

func TestBlogRead(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateMockBlog(t, store)
		got, err := store.Blog.Read(blog.ID)
		if err != nil {
			t.Fatal(err)
		}

		if got.ID != blog.ID {
			t.Fatalf("want %v, got %v", blog.ID, got.ID)
		}

		return test.ErrSkipCommit
	})
}

func TestBlogList(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)

		limit := 3
		offset := 0
		blogs, err := store.Blog.List(limit, offset)
		if err != nil {
			t.Fatal(err)
		}

		if len(blogs) != limit {
			t.Fatalf("want %v, got %v", limit, len(blogs))
		}

		return test.ErrSkipCommit
	})
}
