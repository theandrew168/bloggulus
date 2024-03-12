package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/storage"
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
		blog,
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomTime(),
	)
	return post
}

// mocks a blog and creates it in the database
func CreateMockBlog(t *testing.T, store *storage.Storage) domain.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := store.Blog.Create(blog)
	if err != nil {
		t.Fatal(err)
	}

	return blog
}

// mocks a post and creates it in the database
func CreateMockPost(t *testing.T, store *storage.Storage) domain.Post {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := store.Blog.Create(blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)

	// create an example post
	err = store.Post.Create(post)
	if err != nil {
		t.Fatal(err)
	}

	return post
}
