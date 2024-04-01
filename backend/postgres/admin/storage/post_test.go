package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostCreate(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	// TODO: do something more here?
	store.WithTransaction(func(store storage.Storage) error {
		test.CreateMockPost(t, store)
		return test.ErrRollback
	})
}

func TestPostCreateAlreadyExists(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		post := test.CreateMockPost(t, store)

		err := store.Post().Create(post)
		t.Log(err.Error())
		if !errors.Is(err, postgres.ErrConflict) {
			t.Fatal("duplicate post should return an error")
		}

		return test.ErrRollback
	})
}

func TestPostRead(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		post := test.CreateMockPost(t, store)
		got, err := store.Post().Read(post.ID())
		if err != nil {
			t.Fatal(err)
		}

		if got.ID() != post.ID() {
			t.Fatalf("want %v, got %v", post.ID(), got.ID())
		}

		return test.ErrRollback
	})
}

func TestPostList(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

	store.WithTransaction(func(store storage.Storage) error {
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)

		limit := 3
		offset := 0
		posts, err := store.Post().List(limit, offset)
		if err != nil {
			t.Fatal(err)
		}

		if len(posts) != limit {
			t.Fatalf("want %v, got %v", limit, len(posts))
		}

		return test.ErrRollback
	})
}

func TestPostListByBlog(t *testing.T) {
	store, closer := test.NewAdminStorage(t)
	defer closer()

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
			if err != nil {
				t.Fatal(err)
			}
		}

		limit := 3
		offset := 0
		posts, err := store.Post().ListByBlog(blog, limit, offset)
		if err != nil {
			t.Fatal(err)
		}

		if len(posts) != limit {
			t.Fatalf("want %v, got %v", limit, len(posts))
		}

		// most recent post should be the one just added
		if posts[0].ID() != post.ID() {
			t.Fatalf("want %v, got %v", post.ID(), posts[0].ID())
		}

		return test.ErrRollback
	})
}
