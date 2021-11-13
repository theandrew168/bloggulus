package test

import (
	"context"
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
)

func BlogCreate(t *testing.T, storage core.Storage) {
	// generate some random blog data
	blog := NewMockBlog()
	if blog.BlogID != 0 {
		t.Fatal("blog id before creation should be zero")
	}

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// blog should have an ID after creation
	if blog.BlogID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}

func BlogCreateExists(t *testing.T, storage core.Storage) {
	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// attempt to create the same blog again
	err = storage.BlogCreate(context.Background(), &blog)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate blog should return an error")
	}
}

func BlogReadAll(t *testing.T, storage core.Storage) {
	_, err := storage.BlogReadAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
