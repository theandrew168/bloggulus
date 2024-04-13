package reader_test

import (
	"testing"

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
