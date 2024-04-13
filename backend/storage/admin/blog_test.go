package admin_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewBlog()
		err := store.Admin().Blog().Create(blog)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)

		// attempt to create the same blog again
		err := store.Admin().Blog().Create(blog)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestBlogRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)
		got, err := store.Admin().Blog().Read(blog.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), blog.ID())

		return postgres.ErrRollback
	})
}

func TestBlogReadByFeedURL(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)
		got, err := store.Admin().Blog().ReadByFeedURL(blog.FeedURL())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), blog.ID())

		return postgres.ErrRollback
	})
}

func TestBlogList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)
		test.CreateBlog(t, store)

		limit := 5
		offset := 0
		blogs, err := store.Admin().Blog().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(blogs), limit)

		return postgres.ErrRollback
	})
}

func TestBlogUpdate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)

		etag := "foo"
		blog.SetETag(etag)

		lastModified := "bar"
		blog.SetLastModified(lastModified)

		err := store.Admin().Blog().Update(blog)
		test.AssertNilError(t, err)

		got, err := store.Admin().Blog().Read(blog.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ETag(), etag)
		test.AssertEqual(t, got.LastModified(), lastModified)

		return postgres.ErrRollback
	})
}

func TestBlogDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)

		err := store.Admin().Blog().Delete(blog)
		test.AssertNilError(t, err)

		_, err = store.Admin().Blog().Read(blog.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}