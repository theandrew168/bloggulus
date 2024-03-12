package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostCreate(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	// TODO: do something more here?
	store.WithTransaction(func(store *storage.Storage) error {
		test.CreateMockPost(t, store)
		return test.ErrSkipCommit
	})
}

func TestPostCreateAlreadyExists(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		post := test.CreateMockPost(t, store)

		err := store.Post.Create(post)
		t.Log(err.Error())
		if !errors.Is(err, storage.ErrConflict) {
			t.Fatal("duplicate post should return an error")
		}

		return test.ErrSkipCommit
	})
}

func TestPostRead(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		post := test.CreateMockPost(t, store)
		got, err := store.Post.Read(post.ID)
		if err != nil {
			t.Fatal(err)
		}

		if got.ID != post.ID {
			t.Fatalf("want %v, got %v", post.ID, got.ID)
		}

		return test.ErrSkipCommit
	})
}

func TestPostList(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)
		test.CreateMockPost(t, store)

		limit := 3
		offset := 0
		posts, err := store.Post.List(limit, offset)
		if err != nil {
			t.Fatal(err)
		}

		if len(posts) != limit {
			t.Fatalf("want %v, got %v", limit, len(posts))
		}

		return test.ErrSkipCommit
	})
}

func TestPostListByBlog(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	store.WithTransaction(func(store *storage.Storage) error {
		blog := test.CreateMockBlog(t, store)

		// create 5 posts leaving the most recent one in "post"
		var post domain.Post
		for i := 0; i < 5; i++ {
			post = domain.NewPost(
				blog,
				test.RandomURL(32),
				test.RandomString(32),
				test.RandomString(32),
				test.RandomTime(),
			)
			err := store.Post.Create(post)
			if err != nil {
				t.Fatal(err)
			}
		}

		limit := 3
		offset := 0
		posts, err := store.Post.ListByBlog(blog, limit, offset)
		if err != nil {
			t.Fatal(err)
		}

		if len(posts) != limit {
			t.Fatalf("want %v, got %v", limit, len(posts))
		}

		// most recent post should be the one just added
		if posts[0].ID != post.ID {
			t.Fatalf("want %v, got %v", post.ID, posts[0].ID)
		}

		return test.ErrSkipCommit
	})
}

// func TestPostSearch(t *testing.T) {
// 	store, closer := test.NewStorage(t)
// 	defer closer()

// 	store.WithTransaction(func(store *storage.Storage) error {
// 		blog := test.CreateMockBlog(t, store)
// 		q := "python rust"

// 		// create 5 posts leaving the most recent one in "post"
// 		var post domain.Post
// 		for i := 0; i < 5; i++ {
// 			post = domain.NewPost(
// 				blog,
// 				test.RandomURL(32),
// 				q,
// 				test.RandomString(32),
// 				test.RandomTime(),
// 			)
// 			err := store.Post.Create(post)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		}

// 		limit := 3
// 		offset := 0
// 		posts, err := store.Post.Search(q, limit, offset)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if len(posts) != limit {
// 			t.Fatalf("want %v, got %v", limit, len(posts))
// 		}

// 		// tags will always come back sorted desc
// 		// tags := []string{"Python", "Rust"}
// 		// if !subset(tags, posts[0].Tags) {
// 		// 	t.Fatalf("want superset of %v, got %v", tags, posts[0].Tags)
// 		// }

// 		return test.ErrSkipCommit
// 	})
// }

// func TestPostCount(t *testing.T) {
// 	store, closer := test.NewStorage(t)
// 	defer closer()

// 	store.WithTransaction(func(store *storage.Storage) error {
// 		test.CreateMockPost(t, store)

// 		count, err := store.Post.Count()
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// ensure count is at least one
// 		if count < 1 {
// 			t.Fatalf("want >= 1, got %v", count)
// 		}

// 		return test.ErrSkipCommit
// 	})
// }

// func TestPostCountSearch(t *testing.T) {
// 	store, closer := test.NewStorage(t)
// 	defer closer()

// 	store.WithTransaction(func(store *storage.Storage) error {
// 		// generate some searchable post data
// 		q := "python rust"
// 		blog := test.CreateMockBlog(t, store)
// 		post := domain.NewPost(
// 			test.RandomURL(32),
// 			q,
// 			test.RandomTime(),
// 			test.RandomString(32),
// 			blog,
// 		)

// 		// create a searchable post
// 		err := store.Post.Create(&post)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		count, err := store.Post.CountSearch(q)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// ensure count is at least one
// 		if count < 1 {
// 			t.Fatalf("want >= 1, got %v", count)
// 		}

// 		return test.ErrSkipCommit
// 	})
// }

// check if all items in a exist in b
func subset(a, b []string) bool {
	bset := make(map[string]bool)
	for _, s := range b {
		bset[s] = true
	}

	for _, s := range a {
		if _, ok := bset[s]; !ok {
			return false
		}
	}

	return true
}
