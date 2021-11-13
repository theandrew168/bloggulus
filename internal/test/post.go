package test

import (
	"context"
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
)

func PostCreate(t *testing.T, storage core.Storage) {
	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)
	if post.PostID != 0 {
		t.Fatal("post id before creation should be zero")
	}

	// create an example post
	err = storage.PostCreate(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	if post.PostID == 0 {
		t.Fatal("post id after creation should be nonzero")
	}
}

func PostCreateExists(t *testing.T, storage core.Storage) {
	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)

	// create an example post
	err = storage.PostCreate(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	// attempt to create the same post again
	err = storage.PostCreate(context.Background(), &post)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate post should return an error")
	}
}

func PostReadAllByBlog(t *testing.T, storage core.Storage) {
	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)

	// create an example post
	err = storage.PostCreate(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	posts, err := storage.PostReadAllByBlog(context.Background(), blog.BlogID)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 {
		t.Fatal("expected one post linked to blog")
	}
}

func PostReadRecent(t *testing.T, storage core.Storage) {
	_, err := storage.PostReadRecent(context.Background(), 20, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func PostReadSearch(t *testing.T, storage core.Storage) {
	_, err := storage.PostReadSearch(context.Background(), "", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func PostCountRecent(t *testing.T, storage core.Storage) {
	_, err := storage.PostCountRecent(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func PostCountSearch(t *testing.T, storage core.Storage) {
	_, err := storage.PostCountSearch(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
}
