package suite

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	storageTest "github.com/theandrew168/bloggulus/backend/domain/admin/storage/test"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

func TestBlogCreate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := mock.NewBlog()
		err := store.Blog().Create(blog)
		testutil.AssertNilError(t, err)

		return storage.ErrRollback
	})
}

func TestBlogCreateAlreadyExists(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageTest.CreateMockBlog(t, store)

		// attempt to create the same blog again
		err := store.Blog().Create(blog)
		testutil.AssertErrorIs(t, err, storage.ErrConflict)

		return storage.ErrRollback
	})
}

func TestBlogRead(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageTest.CreateMockBlog(t, store)
		got, err := store.Blog().Read(blog.ID())
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, got.ID(), blog.ID())

		return storage.ErrRollback
	})
}

func TestBlogReadByFeedURL(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageTest.CreateMockBlog(t, store)
		got, err := store.Blog().ReadByFeedURL(blog.FeedURL())
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, got.ID(), blog.ID())

		return storage.ErrRollback
	})
}

func TestBlogList(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		storageTest.CreateMockBlog(t, store)
		storageTest.CreateMockBlog(t, store)
		storageTest.CreateMockBlog(t, store)
		storageTest.CreateMockBlog(t, store)
		storageTest.CreateMockBlog(t, store)

		limit := 3
		offset := 0
		blogs, err := store.Blog().List(limit, offset)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, len(blogs), limit)

		return storage.ErrRollback
	})
}

func TestBlogUpdate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageTest.CreateMockBlog(t, store)

		etag := "foo"
		blog.SetETag(etag)

		lastModified := "bar"
		blog.SetLastModified(lastModified)

		err := store.Blog().Update(blog)
		testutil.AssertNilError(t, err)

		got, err := store.Blog().Read(blog.ID())
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, got.ETag(), etag)
		testutil.AssertEqual(t, got.LastModified(), lastModified)

		return storage.ErrRollback
	})
}

func TestBlogDelete(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageTest.CreateMockBlog(t, store)

		err := store.Blog().Delete(blog)
		testutil.AssertNilError(t, err)

		_, err = store.Blog().Read(blog.ID())
		testutil.AssertErrorIs(t, err, storage.ErrNotFound)

		return storage.ErrRollback
	})
}
