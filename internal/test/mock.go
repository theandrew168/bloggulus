package test

import (
	"testing"

	"github.com/theandrew168/bloggulus"
	"github.com/theandrew168/bloggulus/internal/storage"
)

func NewMockBlog() bloggulus.Blog {
	blog := bloggulus.NewBlog(
		RandomURL(32),
		RandomURL(32),
		RandomString(32),
	)
	return blog
}

func NewMockPost(blog bloggulus.Blog) bloggulus.Post {
	post := bloggulus.NewPost(
		RandomURL(32),
		RandomString(32),
		RandomTime(),
		blog,
	)
	return post
}

// mocks a blog and creates it in the database
func CreateMockBlog(t *testing.T, storage *storage.Storage) bloggulus.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.Blog.Create(&blog)
	if err != nil {
		t.Fatal(err)
	}

	return blog
}

// mocks a post and creates it in the database
func CreateMockPost(t *testing.T, storage *storage.Storage) bloggulus.Post {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.Blog.Create(&blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)

	// create an example post
	err = storage.Post.Create(&post)
	if err != nil {
		t.Fatal(err)
	}

	return post
}
