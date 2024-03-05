package feed

import (
	"io"

	"github.com/theandrew168/bloggulus/internal/domain"
)

type mockReader struct {
	blog  domain.Blog
	posts []domain.Post
}

func NewMockReader(blog domain.Blog, posts []domain.Post) Reader {
	r := mockReader{
		blog:  blog,
		posts: posts,
	}
	return &r
}

func (r *mockReader) ReadBlog(feedURL string) (domain.Blog, error) {
	return r.blog, nil
}

func (r *mockReader) ReadBlogPosts(blog domain.Blog, body io.Reader) ([]domain.Post, error) {
	return r.posts, nil
}

func (r *mockReader) ReadPostBody(post domain.Post) (string, error) {
	return post.Body, nil
}
