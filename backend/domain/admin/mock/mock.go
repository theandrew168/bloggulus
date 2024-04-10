package mock

import (
	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/testutil"
)

func NewBlog() *admin.Blog {
	blog := admin.NewBlog(
		testutil.RandomURL(32),
		testutil.RandomURL(32),
		testutil.RandomString(32),
		testutil.RandomString(32),
		testutil.RandomString(32),
		testutil.RandomTime(),
	)
	return blog
}

func NewPost(blog *admin.Blog) *admin.Post {
	post := admin.NewPost(
		blog,
		testutil.RandomURL(32),
		testutil.RandomString(32),
		testutil.RandomString(32),
		testutil.RandomTime(),
	)
	return post
}

func NewTag() *admin.Tag {
	tag := admin.NewTag(
		testutil.RandomString(32),
	)
	return tag
}
