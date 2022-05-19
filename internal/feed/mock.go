package feed

import (
	"github.com/theandrew168/bloggulus"
)

type mockReader struct {
	blog  bloggulus.Blog
	posts []bloggulus.Post
	body  string
}

func NewMockReader(blog bloggulus.Blog, posts []bloggulus.Post, body string) Reader {
	r := mockReader{
		blog:  blog,
		posts: posts,
		body:  body,
	}
	return &r
}

func (r *mockReader) ReadBlog(feedURL string) (bloggulus.Blog, error) {
	return r.blog, nil
}

func (r *mockReader) ReadBlogPosts(blog bloggulus.Blog) ([]bloggulus.Post, error) {
	return r.posts, nil
}

func (r *mockReader) ReadPostBody(post bloggulus.Post) (string, error) {
	return r.body, nil
}
