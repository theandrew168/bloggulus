package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostCreate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.NewMockBlog()
		err := store.Blog().Create(blog)
		test.AssertNilError(t, err)

		post := test.NewMockPost(blog)
		err = store.Post().Create(post)
		test.AssertNilError(t, err)

		return test.ErrRollback
	})
}

func TestPostCreateAlreadyExists(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := test.CreateMockPost(t, store)

		// attempt to create the same post again
		err := store.Post().Create(post)
		test.AssertErrorIs(t, err, storage.ErrConflict)

		return test.ErrRollback
	})
}

func TestPostRead(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := test.CreateMockPost(t, store)
		got, err := store.Post().Read(post.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), post.ID())

		return test.ErrRollback
	})
}

func TestPostReadByURL(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := test.CreateMockPost(t, store)
		got, err := store.Post().ReadByURL(post.URL())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), post.ID())

		return test.ErrRollback
	})
}

func TestPostList(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)

		limit := 3
		offset := 0
		posts, err := store.Post().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(posts), limit)

		return test.ErrRollback
	})
}

func TestPostListByBlog(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := test.CreateMockBlog(t, store)

		// create 5 posts leaving the most recent one in "post"
		var post *admin.Post
		for i := 0; i < 5; i++ {
			post = admin.NewPost(
				blog,
				test.RandomURL(32),
				test.RandomString(32),
				test.RandomString(32),
				test.RandomTime(),
			)
			err := store.Post().Create(post)
			test.AssertNilError(t, err)
		}

		limit := 3
		offset := 0
		posts, err := store.Post().ListByBlog(blog, limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(posts), limit)

		// most recent post should be the one just added
		test.AssertEqual(t, posts[0].ID(), post.ID())

		return test.ErrRollback
	})
}

func TestPostUpdate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := test.CreateMockPost(t, store)

		contents := "foobar"
		post.SetContents(contents)

		err := store.Post().Update(post)
		test.AssertNilError(t, err)

		got, err := store.Post().Read(post.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.Contents(), contents)

		return test.ErrRollback
	})
}

func TestPostDelete(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := test.CreateMockPost(t, store)

		err := store.Post().Delete(post)
		test.AssertNilError(t, err)

		_, err = store.Post().Read(post.ID())
		test.AssertErrorIs(t, err, storage.ErrNotFound)

		return test.ErrRollback
	})
}
