package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestArticleList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewBlog(t)
		err := store.Blog().Create(blog)
		test.AssertNilError(t, err)

		post := test.NewPost(t, blog)
		err = store.Post().Create(post)
		test.AssertNilError(t, err)

		articles, err := store.Article().List(20, 0)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(articles), 1)
		test.AssertEqual(t, articles[0].Title(), post.Title())
		test.AssertEqual(t, articles[0].URL(), post.URL())
		test.AssertEqual(t, articles[0].BlogTitle(), blog.Title())
		test.AssertEqual(t, articles[0].BlogURL(), blog.SiteURL())

		return postgres.ErrRollback
	})
}

func TestArticleListSearch(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewBlog(t)
		err := store.Blog().Create(blog)
		test.AssertNilError(t, err)

		// create a post about python
		pythonPost, err := model.NewPost(
			blog,
			"https://example.com/python",
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = store.Post().Create(pythonPost)
		test.AssertNilError(t, err)

		// create a post about python
		boringPost, err := model.NewPost(
			blog,
			"https://example.com/boring",
			"Boring",
			"content about nothing",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = store.Post().Create(boringPost)
		test.AssertNilError(t, err)

		// list articles that relate to python
		articles, err := store.Article().ListSearch("python", 20, 0)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(articles), 1)
		test.AssertEqual(t, articles[0].Title(), pythonPost.Title())
		test.AssertEqual(t, articles[0].URL(), pythonPost.URL())

		return postgres.ErrRollback
	})
}

func TestArticleCount(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateBlog(t, store)
		test.CreatePost(t, store, blog)
		test.CreatePost(t, store, blog)
		test.CreatePost(t, store, blog)

		count, err := store.Article().Count()
		test.AssertNilError(t, err)

		test.AssertEqual(t, count, 3)

		return postgres.ErrRollback
	})
}

func TestArticleCountSearch(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.NewBlog(t)
		err := store.Blog().Create(blog)
		test.AssertNilError(t, err)

		// create a post about python
		pythonPost, err := model.NewPost(
			blog,
			"https://example.com/python",
			"Python",
			"content about python",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = store.Post().Create(pythonPost)
		test.AssertNilError(t, err)

		// create a post about python
		boringPost, err := model.NewPost(
			blog,
			"https://example.com/boring",
			"Boring",
			"content about nothing",
			timeutil.Now(),
		)
		test.AssertNilError(t, err)

		err = store.Post().Create(boringPost)
		test.AssertNilError(t, err)

		// count posts that relate to python
		count, err := store.Article().CountSearch("python")
		test.AssertNilError(t, err)

		// should only find one
		test.AssertEqual(t, count, 1)

		return postgres.ErrRollback
	})
}
