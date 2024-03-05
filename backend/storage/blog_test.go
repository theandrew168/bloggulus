package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, storage)
	if blog.ID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, storage)

	// attempt to create the same blog again
	err := storage.Blog.Create(&blog)
	if !errors.Is(err, database.ErrExist) {
		t.Fatal("duplicate blog should return an error")
	}
}

func TestBlogRead(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, storage)
	got, err := storage.Blog.Read(blog.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != blog.ID {
		t.Fatalf("want %v, got %v", blog.ID, got.ID)
	}
}

func TestBlogReadAll(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	test.CreateMockBlog(t, storage)
	test.CreateMockBlog(t, storage)
	test.CreateMockBlog(t, storage)
	test.CreateMockBlog(t, storage)
	test.CreateMockBlog(t, storage)

	limit := 3
	offset := 0
	blogs, err := storage.Blog.ReadAll(limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(blogs) != limit {
		t.Fatalf("want %v, got %v", limit, len(blogs))
	}
}
