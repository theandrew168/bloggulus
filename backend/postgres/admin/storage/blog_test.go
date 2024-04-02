package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.NewMockBlog()
		err := store.Blog().Create(blog)
		if err != nil {
			t.Fatal(err)
		}

		return test.ErrRollback
	})
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)

		// attempt to create the same blog again
		err := store.Blog().Create(blog)
		if !errors.Is(err, postgres.ErrConflict) {
			t.Fatal("duplicate blog should return an error")
		}

		return test.ErrRollback
	})
}

func TestBlogRead(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)
		got, err := store.Blog().Read(blog.ID())
		if err != nil {
			t.Fatal(err)
		}

		if got.ID() != blog.ID() {
			t.Fatalf("want %v, got %v", blog.ID(), got.ID())
		}

		return test.ErrRollback
	})
}

func TestBlogReadByFeedURL(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)
		got, err := store.Blog().ReadByFeedURL(blog.FeedURL())
		if err != nil {
			t.Fatal(err)
		}

		if got.ID() != blog.ID() {
			t.Fatalf("want %v, got %v", blog.ID(), got.ID())
		}

		return test.ErrRollback
	})
}

func TestBlogList(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)
		test.CreateMockBlog(t, store)

		limit := 3
		offset := 0
		blogs, err := store.Blog().List(limit, offset)
		if err != nil {
			t.Fatal(err)
		}

		if len(blogs) != limit {
			t.Fatalf("want %v, got %v", limit, len(blogs))
		}

		return test.ErrRollback
	})
}

func TestBlogUpdate(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.NewMockBlog()
		err := store.Blog().Create(blog)
		if err != nil {
			t.Fatal(err)
		}

		etag := "foo"
		blog.SetETag(etag)

		lastModified := "bar"
		blog.SetLastModified(lastModified)

		err = store.Blog().Update(blog)
		if err != nil {
			t.Fatal(err)
		}

		got, err := store.Blog().Read(blog.ID())
		if err != nil {
			t.Fatal(err)
		}

		if got.ETag() != etag {
			t.Fatalf("want %v, got %v", lastModified, got.ETag())
		}

		if got.LastModified() != lastModified {
			t.Fatalf("want %v, got %v", lastModified, got.LastModified())
		}

		return test.ErrRollback
	})
}
