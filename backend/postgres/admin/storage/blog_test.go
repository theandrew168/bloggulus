package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.NewMockBlog()
		err := store.Blog().Create(blog)
		test.AssertNilError(t, err)

		return test.ErrRollback
	})
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)

		// attempt to create the same blog again
		err := store.Blog().Create(blog)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return test.ErrRollback
	})
}

func TestBlogRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)
		got, err := store.Blog().Read(blog.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), blog.ID())

		return test.ErrRollback
	})
}

func TestBlogReadByFeedURL(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)
		got, err := store.Blog().ReadByFeedURL(blog.FeedURL())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), blog.ID())

		return test.ErrRollback
	})
}

func TestBlogList(t *testing.T) {
	t.Parallel()

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
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(blogs), limit)

		return test.ErrRollback
	})
}

func TestBlogUpdate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)

		etag := "foo"
		blog.SetETag(etag)

		lastModified := "bar"
		blog.SetLastModified(lastModified)

		err := store.Blog().Update(blog)
		test.AssertNilError(t, err)

		got, err := store.Blog().Read(blog.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ETag(), etag)
		test.AssertEqual(t, got.LastModified(), lastModified)

		return test.ErrRollback
	})
}

func TestBlogDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)

		err := store.Blog().Delete(blog)
		test.AssertNilError(t, err)

		_, err = store.Blog().Read(blog.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return test.ErrRollback
	})
}
