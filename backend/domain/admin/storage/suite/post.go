package suite

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	storageTest "github.com/theandrew168/bloggulus/backend/domain/admin/storage/test"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostCreate(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := mock.NewBlog()
		err := store.Blog().Create(blog)
		test.AssertNilError(t, err)

		post := mock.NewPost(blog)
		err = store.Post().Create(post)
		test.AssertNilError(t, err)

		return storage.ErrRollback
	})
}

func TestPostCreateAlreadyExists(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := storageTest.CreateMockPost(t, store)

		// attempt to create the same post again
		err := store.Post().Create(post)
		test.AssertErrorIs(t, err, storage.ErrConflict)

		return storage.ErrRollback
	})
}

func TestPostRead(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := storageTest.CreateMockPost(t, store)
		got, err := store.Post().Read(post.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), post.ID())

		return storage.ErrRollback
	})
}

func TestPostReadByURL(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := storageTest.CreateMockPost(t, store)
		got, err := store.Post().ReadByURL(post.URL())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.ID(), post.ID())

		return storage.ErrRollback
	})
}

func TestPostList(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		storageTest.CreateMockPost(t, store)
		storageTest.CreateMockPost(t, store)
		storageTest.CreateMockPost(t, store)
		storageTest.CreateMockPost(t, store)
		storageTest.CreateMockPost(t, store)

		limit := 5
		offset := 0
		posts, err := store.Post().List(limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(posts), limit)

		return storage.ErrRollback
	})
}

func TestPostListByBlog(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		blog := storageTest.CreateMockBlog(t, store)

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

		limit := 5
		offset := 0
		posts, err := store.Post().ListByBlog(blog, limit, offset)
		test.AssertNilError(t, err)

		test.AssertEqual(t, len(posts), limit)

		// most recent post should be the one just added
		test.AssertEqual(t, posts[0].ID(), post.ID())

		return storage.ErrRollback
	})
}

func TestPostUpdate(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := storageTest.CreateMockPost(t, store)

		contents := "foobar"
		post.SetContents(contents)

		err := store.Post().Update(post)
		test.AssertNilError(t, err)

		got, err := store.Post().Read(post.ID())
		test.AssertNilError(t, err)

		test.AssertEqual(t, got.Contents(), contents)

		return storage.ErrRollback
	})
}

func TestPostDelete(t *testing.T, store storage.Storage) {
	t.Parallel()

	store.WithTransaction(func(store storage.Storage) error {
		post := storageTest.CreateMockPost(t, store)

		err := store.Post().Delete(post)
		test.AssertNilError(t, err)

		_, err = store.Post().Read(post.ID())
		test.AssertErrorIs(t, err, storage.ErrNotFound)

		return storage.ErrRollback
	})
}
