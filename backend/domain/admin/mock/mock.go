package mock

import (
	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/test"
)

func NewBlog() *admin.Blog {
	blog := admin.NewBlog(
		test.RandomURL(32),
		test.RandomURL(32),
		test.RandomString(32),
		test.RandomString(32),
		test.RandomString(32),
		test.RandomTime(),
	)
	return blog
}

func NewPost(blog *admin.Blog) *admin.Post {
	post := admin.NewPost(
		blog,
		test.RandomURL(32),
		test.RandomString(32),
		test.RandomString(32),
		test.RandomTime(),
	)
	return post
}

func NewTag() *admin.Tag {
	tag := admin.NewTag(
		test.RandomString(32),
	)
	return tag
}
