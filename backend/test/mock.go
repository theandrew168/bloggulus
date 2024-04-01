package test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	adminStorage "github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

func NewMockBlog() admin.Blog {
	blog := admin.NewBlog(
		RandomURL(32),
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomString(32),
		RandomTime(),
	)
	return blog
}

func NewMockPost(blog admin.Blog) admin.Post {
	post := admin.NewPost(
		blog,
		RandomURL(32),
		RandomString(32),
		RandomString(32),
		RandomTime(),
	)
	return post
}

func NewMockTag() *admin.Tag {
	tag := admin.NewTag(
		RandomString(32),
	)
	return tag
}

// mocks a blog and creates it in the database
func CreateMockBlog(t *testing.T, store adminStorage.Storage) admin.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := store.Blog().Create(blog)
	if err != nil {
		t.Fatal(err)
	}

	return blog
}

// mocks a post and creates it in the database
func CreateMockPost(t *testing.T, store adminStorage.Storage) admin.Post {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := store.Blog().Create(blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)

	// create an example post
	err = store.Post().Create(post)
	if err != nil {
		t.Fatal(err)
	}

	return post
}

// mocks a tag and creates it in the database
func CreateMockTag(t *testing.T, store adminStorage.Storage) *admin.Tag {
	t.Helper()

	// generate some random tag data
	tag := NewMockTag()

	// create an example blog
	err := store.Tag().Create(tag)
	if err != nil {
		t.Fatal(err)
	}

	return tag
}
