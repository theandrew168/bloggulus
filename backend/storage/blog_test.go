package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

func TestBlogCreate(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, store)
	if blog.ID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}

func TestBlogCreateAlreadyExists(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, store)

	// attempt to create the same blog again
	err := store.Blog.Create(&blog)
	if !errors.Is(err, storage.ErrExist) {
		t.Fatal("duplicate blog should return an error")
	}
}

func TestBlogRead(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, store)
	got, err := store.Blog.Read(blog.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != blog.ID {
		t.Fatalf("want %v, got %v", blog.ID, got.ID)
	}
}

func TestBlogReadAll(t *testing.T) {
	store, closer := test.NewStorage(t)
	defer closer()

	test.CreateMockBlog(t, store)
	test.CreateMockBlog(t, store)
	test.CreateMockBlog(t, store)
	test.CreateMockBlog(t, store)
	test.CreateMockBlog(t, store)

	limit := 3
	offset := 0
	blogs, err := store.Blog.ReadAll(limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(blogs) != limit {
		t.Fatalf("want %v, got %v", limit, len(blogs))
	}
}
