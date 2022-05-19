package storage_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/internal/test"
)

func TestBlogCreate(t *testing.T) {
	storage, closer := test.NewStorage(t)
	defer closer()

	blog := test.CreateMockBlog(t, storage)

	// blog should have an ID after creation
	if blog.ID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}
