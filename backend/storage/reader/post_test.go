package reader_test

import (
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewBlog()
		err := store.Admin().Blog().Create(blog)
		test.AssertNilError(t, err)

		post := test.NewPost(blog)
		err = store.Admin().Post().Create(post)
		test.AssertNilError(t, err)

		posts, err := store.Reader().Post().List(20, 0)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(posts), 1)
		test.AssertEqual(t, posts[0].Title(), post.Title())
		test.AssertEqual(t, posts[0].URL(), post.URL())
		test.AssertEqual(t, posts[0].BlogTitle(), blog.Title())
		test.AssertEqual(t, posts[0].BlogURL(), blog.SiteURL())

		return postgres.ErrRollback
	})
}

func TestPostListSearch(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewBlog()
		err := store.Admin().Blog().Create(blog)
		test.AssertNilError(t, err)

		// create a post about python
		pythonPost := admin.NewPost(
			blog,
			"https://example.com/python",
			"Python",
			"content about python",
			time.Now().UTC(),
		)
		err = store.Admin().Post().Create(pythonPost)
		test.AssertNilError(t, err)

		// create a post about python
		boringPost := admin.NewPost(
			blog,
			"https://example.com/boring",
			"Boring",
			"content about nothing",
			time.Now().UTC(),
		)
		err = store.Admin().Post().Create(boringPost)
		test.AssertNilError(t, err)

		// list posts that relate to python
		posts, err := store.Reader().Post().ListSearch("python", 20, 0)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(posts), 1)
		test.AssertEqual(t, posts[0].Title(), pythonPost.Title())
		test.AssertEqual(t, posts[0].URL(), pythonPost.URL())

		return postgres.ErrRollback
	})
}

func TestPostCount(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		test.CreatePost(t, store)
		test.CreatePost(t, store)
		test.CreatePost(t, store)

		count, err := store.Reader().Post().Count()
		test.AssertNilError(t, err)

		test.AssertEqual(t, count, 3)

		return postgres.ErrRollback
	})
}

func TestPostCountSearch(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewBlog()
		err := store.Admin().Blog().Create(blog)
		test.AssertNilError(t, err)

		// create a post about python
		pythonPost := admin.NewPost(
			blog,
			"https://example.com/python",
			"Python",
			"content about python",
			time.Now().UTC(),
		)
		err = store.Admin().Post().Create(pythonPost)
		test.AssertNilError(t, err)

		// create a post about python
		boringPost := admin.NewPost(
			blog,
			"https://example.com/boring",
			"Boring",
			"content about nothing",
			time.Now().UTC(),
		)
		err = store.Admin().Post().Create(boringPost)
		test.AssertNilError(t, err)

		// count posts that relate to python
		count, err := store.Reader().Post().CountSearch("python")
		test.AssertNilError(t, err)

		// should only find one
		test.AssertEqual(t, count, 1)

		return postgres.ErrRollback
	})
}
