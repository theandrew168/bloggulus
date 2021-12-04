package test

import (
	"context"
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
)

func CreateBlog(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)

	// blog should have an ID after creation
	if blog.ID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}

func CreateBlogAlreadyExists(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)

	// attempt to create the same blog again
	err := storage.CreateBlog(context.Background(), &blog)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate blog should return an error")
	}
}

func ReadBlog(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)

	got, err := storage.ReadBlog(context.Background(), blog.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != blog.ID {
		t.Fatalf("want %v, got %v", blog.ID, got.ID)
	}
}

func ReadBlogs(storage core.Storage, t *testing.T) {
	CreateMockBlog(storage, t)
	CreateMockBlog(storage, t)
	CreateMockBlog(storage, t)
	CreateMockBlog(storage, t)
	CreateMockBlog(storage, t)

	limit := 3
	offset := 0
	blogs, err := storage.ReadBlogs(context.Background(), limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(blogs) != limit {
		t.Fatalf("want %v, got %v", limit, len(blogs))
	}
}

func CreateMockBlog(storage core.Storage, t *testing.T) core.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.CreateBlog(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	return blog
}
