package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/internal/domain"
	"github.com/theandrew168/bloggulus/internal/storage"
)

func NewMockBlog() domain.Blog {
	blog := domain.NewBlog(
		RandomURL(32),
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomString(32),
	)
	return blog
}

func NewMockPost(blog domain.Blog) domain.Post {
	post := domain.NewPost(
		RandomURL(32),
		RandomString(32),
		RandomTime(),
		RandomString(32),
		blog,
	)
	return post
}

// mocks a blog and creates it in the database
func CreateMockBlog(t *testing.T, storage *storage.Storage) domain.Blog {
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
func CreateMockPost(t *testing.T, storage *storage.Storage) domain.Post {
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
