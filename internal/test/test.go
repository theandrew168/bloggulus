package test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/feed"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return string(buf)
}

func RandomURL(n int) string {
	return "https://" + RandomString(n)
}

func RandomTime() time.Time {
	return time.Now()
}

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

func ConnectDB(t *testing.T) *pgxpool.Pool {
	// check for database connection url var
	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}

	return conn
}

type mockReader struct {
	blog  core.Blog
	posts []core.Post
	body  string
}

func NewMockFeedReader(blog core.Blog, posts []core.Post, body string) feed.Reader {
	r := mockReader{
		blog:  blog,
		posts: posts,
		body:  body,
	}
	return &r
}

func (r *mockReader) ReadBlog(feedURL string) (core.Blog, error) {
	return r.blog, nil
}

func (r *mockReader) ReadBlogPosts(blog core.Blog) ([]core.Post, error) {
	return r.posts, nil
}

func (r *mockReader) ReadPostBody(post core.Post) (string, error) {
	return r.body, nil
}
