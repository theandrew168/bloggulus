package feed

import (
	"io"

	"github.com/theandrew168/bloggulus"
)

type mockReader struct {
	blog  bloggulus.Blog
	posts []bloggulus.Post
}

func NewMockReader(blog bloggulus.Blog, posts []bloggulus.Post) Reader {
	r := mockReader{
		blog:  blog,
		posts: posts,
	}
	return &r
}

func (r *mockReader) ReadBlog(feedURL string) (bloggulus.Blog, error) {
	return r.blog, nil
}

func (r *mockReader) ReadBlogPosts(blog bloggulus.Blog, body io.Reader) ([]bloggulus.Post, error) {
	return r.posts, nil
}

func (r *mockReader) ReadPostBody(post bloggulus.Post) (string, error) {
	return post.Body, nil
}
