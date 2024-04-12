package postgres_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	storageMock "github.com/theandrew168/bloggulus/backend/domain/admin/storage/mock"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := mock.NewBlog()
		err := store.Blog().Create(blog)
		test.AssertNilError(t, err)

		return postgres.ErrRollback
	})
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageMock.CreateBlog(t, store)

		// attempt to create the same blog again
		err := store.Blog().Create(blog)
		test.AssertErrorIs(t, err, postgres.ErrConflict)

		return postgres.ErrRollback
	})
}

func TestBlogRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageMock.CreateBlog(t, store)
		got, err := store.Blog().Read(blog.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), blog.ID())

		return postgres.ErrRollback
	})
}

func TestBlogReadByFeedURL(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageMock.CreateBlog(t, store)
		got, err := store.Blog().ReadByFeedURL(blog.FeedURL())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), blog.ID())

		return postgres.ErrRollback
	})
}

func TestBlogList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		storageMock.CreateBlog(t, store)
		storageMock.CreateBlog(t, store)
		storageMock.CreateBlog(t, store)
		storageMock.CreateBlog(t, store)
		storageMock.CreateBlog(t, store)

		limit := 5
		offset := 0
		blogs, err := store.Blog().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(blogs), limit)

		return postgres.ErrRollback
	})
}

func TestBlogUpdate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageMock.CreateBlog(t, store)

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

		return postgres.ErrRollback
	})
}

func TestBlogDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageMock.CreateBlog(t, store)

		err := store.Blog().Delete(blog)
		test.AssertNilError(t, err)

		_, err = store.Blog().Read(blog.ID())
		test.AssertErrorIs(t, err, postgres.ErrNotFound)

		return postgres.ErrRollback
	})
}
