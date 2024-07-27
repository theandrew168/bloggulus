package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.NewBlog(t)
	err := store.Blog().Create(blog)
	test.AssertNilError(t, err)
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)

	// attempt to create the same blog again
	err := store.Blog().Create(blog)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestBlogRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	got, err := store.Blog().Read(blog.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), blog.ID())
}

func TestBlogReadByFeedURL(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	got, err := store.Blog().ReadByFeedURL(blog.FeedURL())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), blog.ID())
}

func TestBlogList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	test.CreateBlog(t, store)
	test.CreateBlog(t, store)
	test.CreateBlog(t, store)

	limit := 3
	offset := 0
	blogs, err := store.Blog().List(limit, offset)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(blogs), limit)
}

func TestBlogListAll(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	test.CreateBlog(t, store)
	test.CreateBlog(t, store)
	test.CreateBlog(t, store)

	blogs, err := store.Blog().ListAll()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, len(blogs), 3)
}

func TestBlogCount(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	test.CreateBlog(t, store)
	test.CreateBlog(t, store)
	test.CreateBlog(t, store)

	count, err := store.Blog().Count()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, count, 3)
}

func TestBlogUpdate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)

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
}

func TestBlogDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)

	err := store.Blog().Delete(blog)
	test.AssertNilError(t, err)

	_, err = store.Blog().Read(blog.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
