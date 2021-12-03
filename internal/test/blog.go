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
		t.Error("blog id after creation should be nonzero")
	}
}

func CreateBlogAlreadyExists(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)

	// attempt to create the same blog again
	err := storage.CreateBlog(context.Background(), &blog)
	if !errors.Is(err, core.ErrExist) {
		t.Error("duplicate blog should return an error")
	}
}

func ReadBlog(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)

	got, err := storage.ReadBlog(context.Background(), blog.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != blog.ID {
		t.Errorf("want %v, got %v", blog.ID, got.ID)
	}
}

// TODO: test pagination
func ReadBlogs(storage core.Storage, t *testing.T) {
	CreateMockBlog(storage, t)

	blogs, err := storage.ReadBlogs(context.Background(), 20, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(blogs) < 1 {
		t.Errorf("want >= 1, got %v", len(blogs))
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
