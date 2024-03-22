package feed

import (
	"io"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
)

type mockReader struct {
	blog  admin.Blog
	posts []admin.Post
}

func NewMockReader(blog admin.Blog, posts []admin.Post) Reader {
	r := mockReader{
		blog:  blog,
		posts: posts,
	}
	return &r
}

func (r *mockReader) ReadBlog(feedURL string) (admin.Blog, error) {
	return r.blog, nil
}

func (r *mockReader) ReadBlogPosts(blog admin.Blog, body io.Reader) ([]admin.Post, error) {
	return r.posts, nil
}

func (r *mockReader) ReadPostBody(post admin.Post) (string, error) {
	return post.Content, nil
}
