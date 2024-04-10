package todo

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/mock"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

// mocks a blog and creates it in the database
func CreateMockBlog(t *testing.T, store storage.Storage) *admin.Blog {
	t.Helper()

	// generate some random blog data
	blog := mock.NewBlog()

	// create an example blog
	err := store.Blog().Create(blog)
	testutil.AssertNilError(t, err)

	return blog
}

// mocks a post and creates it in the database
func CreateMockPost(t *testing.T, store storage.Storage) *admin.Post {
	t.Helper()

	// generate some random blog data
	blog := mock.NewBlog()

	// create an example blog
	err := store.Blog().Create(blog)
	testutil.AssertNilError(t, err)

	// generate some random post data
	post := mock.NewPost(blog)

	// create an example post
	err = store.Post().Create(post)
	testutil.AssertNilError(t, err)

	return post
}

// mocks a tag and creates it in the database
func CreateMockTag(t *testing.T, store storage.Storage) *admin.Tag {
	t.Helper()

	// generate some random tag data
	tag := mock.NewTag()

	// create an example blog
	err := store.Tag().Create(tag)
	testutil.AssertNilError(t, err)

	return tag
}
