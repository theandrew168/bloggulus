package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

func TestArticleList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	test.CreatePost(t, store, blog)

	articles, err := store.Article().List(1, 0)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(articles), 1)
}

func TestArticleListSearch(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.NewBlog(t)
	err := store.Blog().Create(blog)
	test.AssertNilError(t, err)

	// create a post about python
	pythonPost, err := model.NewPost(
		blog,
		test.RandomURL(20),
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
		test.RandomURL(20),
		"Boring",
		"content about nothing",
		timeutil.Now(),
	)
	test.AssertNilError(t, err)

	err = store.Post().Create(boringPost)
	test.AssertNilError(t, err)

	// list articles that relate to python
	articles, err := store.Article().ListSearch("python", 1, 0)
	test.AssertNilError(t, err)

	// should find at least one
	test.AssertEqual(t, len(articles), 1)
}

func TestArticleCount(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	test.CreatePost(t, store, blog)
	test.CreatePost(t, store, blog)
	test.CreatePost(t, store, blog)

	count, err := store.Article().Count()
	test.AssertNilError(t, err)

	test.AssertAtLeast(t, count, 3)
}

func TestArticleCountSearch(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.NewBlog(t)
	err := store.Blog().Create(blog)
	test.AssertNilError(t, err)

	// create a post about python
	pythonPost, err := model.NewPost(
		blog,
		test.RandomURL(20),
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
		test.RandomURL(20),
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

	// should find at least one
	test.AssertAtLeast(t, count, 1)
}
