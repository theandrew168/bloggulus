package suite

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

func RunStorageTests(t *testing.T, store storage.Storage) {
	tests := []struct {
		name string
		test func(t *testing.T, store storage.Storage)
	}{
		{name: "TestBlogCreate", test: TestBlogCreate},
		{name: "TestBlogCreateAlreadyExists", test: TestBlogCreateAlreadyExists},
		{name: "TestBlogRead", test: TestBlogRead},
		{name: "TestBlogReadByFeedURL", test: TestBlogReadByFeedURL},
		{name: "TestBlogList", test: TestBlogList},
		{name: "TestBlogUpdate", test: TestBlogUpdate},
		{name: "TestBlogDelete", test: TestBlogDelete},

		{name: "TestPostCreate", test: TestPostCreate},
		{name: "TestPostCreateAlreadyExists", test: TestPostCreateAlreadyExists},
		{name: "TestPostRead", test: TestPostRead},
		{name: "TestPostReadByURL", test: TestPostReadByURL},
		{name: "TestPostList", test: TestPostList},
		{name: "TestPostListByBlog", test: TestPostListByBlog},
		{name: "TestPostUpdate", test: TestPostUpdate},
		{name: "TestPostDelete", test: TestPostDelete},

		{name: "TestTagCreate", test: TestTagCreate},
		{name: "TestTagCreateAlreadyExists", test: TestTagCreateAlreadyExists},
		{name: "TestTagRead", test: TestTagRead},
		{name: "TestTagList", test: TestTagList},
		{name: "TestTagDelete", test: TestTagDelete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, store)
		})
	}
}
