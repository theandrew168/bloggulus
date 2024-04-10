package suite

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage/todo"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

func TestPostCreate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := mock.NewBlog()
		err := store.Blog().Create(blog)
		testutil.AssertNilError(t, err)

		post := mock.NewPost(blog)
		err = store.Post().Create(post)
		testutil.AssertNilError(t, err)

		return storage.ErrRollback
	})
}

func TestPostCreateAlreadyExists(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := todo.CreateMockPost(t, store)

		// attempt to create the same post again
		err := store.Post().Create(post)
		testutil.AssertErrorIs(t, err, storage.ErrConflict)

		return storage.ErrRollback
	})
}

func TestPostRead(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := todo.CreateMockPost(t, store)
		got, err := store.Post().Read(post.ID())
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, got.ID(), post.ID())

		return storage.ErrRollback
	})
}

func TestPostReadByURL(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := todo.CreateMockPost(t, store)
		got, err := store.Post().ReadByURL(post.URL())
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, got.ID(), post.ID())

		return storage.ErrRollback
	})
}

func TestPostList(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)
		todo.CreateMockPost(t, store)

		limit := 3
		offset := 0
		posts, err := store.Post().List(limit, offset)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, len(posts), limit)

		return storage.ErrRollback
	})
}

func TestPostListByBlog(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := todo.CreateMockBlog(t, store)

		// create 5 posts leaving the most recent one in "post"
		var post *admin.Post
		for i := 0; i < 5; i++ {
			post = admin.NewPost(
				blog,
				testutil.RandomURL(32),
				testutil.RandomString(32),
				testutil.RandomString(32),
				testutil.RandomTime(),
			)
			err := store.Post().Create(post)
			testutil.AssertNilError(t, err)
		}

		limit := 3
		offset := 0
		posts, err := store.Post().ListByBlog(blog, limit, offset)
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, len(posts), limit)

		// most recent post should be the one just added
		testutil.AssertEqual(t, posts[0].ID(), post.ID())

		return storage.ErrRollback
	})
}

func TestPostUpdate(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := todo.CreateMockPost(t, store)

		contents := "foobar"
		post.SetContents(contents)

		err := store.Post().Update(post)
		testutil.AssertNilError(t, err)

		got, err := store.Post().Read(post.ID())
		testutil.AssertNilError(t, err)

		testutil.AssertEqual(t, got.Contents(), contents)

		return storage.ErrRollback
	})
}

func TestPostDelete(t *testing.T, store storage.Storage) {
	// t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := todo.CreateMockPost(t, store)

		err := store.Post().Delete(post)
		testutil.AssertNilError(t, err)

		_, err = store.Post().Read(post.ID())
		testutil.AssertErrorIs(t, err, storage.ErrNotFound)

		return storage.ErrRollback
	})
}
