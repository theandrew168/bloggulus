package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.NewBlog(t)
	err := store.Blog().Create(blog)
	test.AssertNilError(t, err)

	post := test.NewPost(t, blog)
	err = store.Post().Create(post)
	test.AssertNilError(t, err)
}

func TestPostCreateAlreadyExists(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	post := test.CreatePost(t, store, blog)

	// attempt to create the same post again
	err := store.Post().Create(post)
	test.AssertErrorIs(t, err, postgres.ErrConflict)
}

func TestPostRead(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	post := test.CreatePost(t, store, blog)
	got, err := store.Post().Read(post.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), post.ID())
}

func TestPostReadByURL(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	post := test.CreatePost(t, store, blog)
	got, err := store.Post().ReadByURL(post.URL())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID(), post.ID())
}

func TestPostList(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	test.CreatePost(t, store, blog)
	test.CreatePost(t, store, blog)
	test.CreatePost(t, store, blog)

	limit := 3
	offset := 0
	posts, err := store.Post().List(blog, limit, offset)
	test.AssertNilError(t, err)

	test.AssertEqual(t, len(posts), limit)
}

func TestPostCount(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	test.CreatePost(t, store, blog)
	test.CreatePost(t, store, blog)
	test.CreatePost(t, store, blog)

	count, err := store.Post().Count(blog)
	test.AssertNilError(t, err)

	test.AssertEqual(t, count, 3)
}

func TestPostUpdate(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	post := test.CreatePost(t, store, blog)

	content := "foobar"
	post.SetContent(content)

	err := store.Post().Update(post)
	test.AssertNilError(t, err)

	got, err := store.Post().Read(post.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.Content(), content)
}

func TestPostDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateBlog(t, store)
	post := test.CreatePost(t, store, blog)

	err := store.Post().Delete(post)
	test.AssertNilError(t, err)

	_, err = store.Post().Read(post.ID())
	test.AssertErrorIs(t, err, postgres.ErrNotFound)
}
