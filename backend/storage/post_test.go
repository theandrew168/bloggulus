package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestPostCreate(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	post := test.CreateMockPost(t, storage)
	if post.ID == 0 {
		t.Fatal("post id after creation should be nonzero")
	}
}

func TestPostCreateAlreadyExists(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	post := test.CreateMockPost(t, storage)

	err := storage.Post.Create(&post)
	if !errors.Is(err, database.ErrExist) {
		t.Fatal("duplicate post should return an error")
	}
}

func TestPostRead(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	post := test.CreateMockPost(t, storage)
	got, err := storage.Post.Read(post.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != post.ID {
		t.Fatalf("want %v, got %v", post.ID, got.ID)
	}
}

func TestPostReadAll(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)
	test.CreateMockPost(t, storage)

	limit := 3
	offset := 0
	posts, err := storage.Post.ReadAll(limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != limit {
		t.Fatalf("want %v, got %v", limit, len(posts))
	}
}

func TestPostReadAllByBlog(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, storage)

	// create 5 posts leaving the most recent one in "post"
	var post domain.Post
	for i := 0; i < 5; i++ {
		post = domain.NewPost(
			test.RandomURL(32),
			test.RandomString(32),
			test.RandomTime(),
			test.RandomString(32),
			blog,
		)
		err := storage.Post.Create(&post)
		if err != nil {
			t.Fatal(err)
		}
	}

	limit := 3
	offset := 0
	posts, err := storage.Post.ReadAllByBlog(blog, limit, offset)
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
}

func TestPostSearch(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, storage)
	q := "python rust"

	// create 5 posts leaving the most recent one in "post"
	var post domain.Post
	for i := 0; i < 5; i++ {
		post = domain.NewPost(
			test.RandomURL(32),
			q,
			test.RandomTime(),
			test.RandomString(32),
			blog,
		)
		err := storage.Post.Create(&post)
		if err != nil {
			t.Fatal(err)
		}
	}

	limit := 3
	offset := 0
	posts, err := storage.Post.Search(q, limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != limit {
		t.Fatalf("want %v, got %v", limit, len(posts))
	}

	// tags will always come back sorted desc
	tags := []string{"Python", "Rust"}
	if !subset(tags, posts[0].Tags) {
		t.Fatalf("want superset of %v, got %v", tags, posts[0].Tags)
	}
}

func TestPostCount(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	test.CreateMockPost(t, storage)

	count, err := storage.Post.Count()
	if err != nil {
		t.Fatal(err)
	}

	// ensure count is at least one
	if count < 1 {
		t.Fatalf("want >= 1, got %v", count)
	}
}

func TestPostCountSearch(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	// generate some searchable post data
	q := "python rust"
	blog := test.CreateMockBlog(t, storage)
	post := domain.NewPost(
		test.RandomURL(32),
		q,
		test.RandomTime(),
		test.RandomString(32),
		blog,
	)

	// create a searchable post
	err := storage.Post.Create(&post)
	if err != nil {
		t.Fatal(err)
	}

	count, err := storage.Post.CountSearch(q)
	if err != nil {
		t.Fatal(err)
	}

	// ensure count is at least one
	if count < 1 {
		t.Fatalf("want >= 1, got %v", count)
	}
}

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
