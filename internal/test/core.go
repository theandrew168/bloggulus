package test

import (
	"github.com/theandrew168/bloggulus/internal/core"
)

func NewMockBlog() core.Blog {
	blog := core.NewBlog(
		RandomURL(32),
		RandomURL(32),
		RandomString(32),
	)
	return blog
}

func NewMockPost(blog core.Blog) core.Post {
	post := core.NewPost(
		RandomURL(32),
		RandomString(32),
		RandomTime(),
		blog,
	)
	return post
}
