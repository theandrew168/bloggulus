package mock

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/test"
)

// mocks a blog and creates it in the database
func CreateBlog(t *testing.T, store storage.Storage) *admin.Blog {
	t.Helper()

	// generate some random blog data
	blog := mock.NewBlog()

	// create an example blog
	err := store.Blog().Create(blog)
	test.AssertNilError(t, err)

	return blog
}

// mocks a post and creates it in the database
func CreatePost(t *testing.T, store storage.Storage) *admin.Post {
	t.Helper()

	// generate some random blog data
	blog := mock.NewBlog()

	// create an example blog
	err := store.Blog().Create(blog)
	test.AssertNilError(t, err)

	// generate some random post data
	post := mock.NewPost(blog)

	// create an example post
	err = store.Post().Create(post)
	test.AssertNilError(t, err)

	return post
}

// mocks a tag and creates it in the database
func CreateTag(t *testing.T, store storage.Storage) *admin.Tag {
	t.Helper()

	// generate some random tag data
	tag := mock.NewTag()

	// create an example blog
	err := store.Tag().Create(tag)
	test.AssertNilError(t, err)

	return tag
}