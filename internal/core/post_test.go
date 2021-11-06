package core_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/postgresql"
)

func TestPostCreate(t *testing.T) {
	conn := connectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)
	postStorage := postgresql.NewPostStorage(conn)

	// generate some random blog data
	feedURL := "https://" + randomString(32)
	siteURL := "https://" + randomString(32)
	blogTitle := randomString(32)

	blog := core.NewBlog(feedURL, siteURL, blogTitle)

	// create an example blog
	err := blogStorage.Create(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	url := "https://" + randomString(32)
	title := randomString(32)
	updated := time.Now()

	post := core.NewPost(url, title, updated, blog)
	if post.PostID != 0 {
		t.Fatal("post id before creation should be zero")
	}

	// create an example post
	err = postStorage.Create(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	if post.PostID == 0 {
		t.Fatal("post id after creation should be nonzero")
	}
}

func TestPostCreateExists(t *testing.T) {
	conn := connectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)
	postStorage := postgresql.NewPostStorage(conn)

	// generate some random blog data
	feedURL := "https://" + randomString(32)
	siteURL := "https://" + randomString(32)
	blogTitle := randomString(32)

	blog := core.NewBlog(feedURL, siteURL, blogTitle)

	// create an example blog
	err := blogStorage.Create(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	url := "https://" + randomString(32)
	title := randomString(32)
	updated := time.Now()

	post := core.NewPost(url, title, updated, blog)

	// create an example post
	err = postStorage.Create(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	// attempt to create the same post again
	err = postStorage.Create(context.Background(), &post)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate post should return an error")
	}
}

func TestPostReadAllByBlog(t *testing.T) {
	conn := connectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)
	postStorage := postgresql.NewPostStorage(conn)

	// generate some random blog data
	feedURL := "https://" + randomString(32)
	siteURL := "https://" + randomString(32)
	blogTitle := randomString(32)

	blog := core.NewBlog(feedURL, siteURL, blogTitle)

	// create an example blog
	err := blogStorage.Create(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	url := "https://" + randomString(32)
	title := randomString(32)
	updated := time.Now()

	post := core.NewPost(url, title, updated, blog)

	// create an example post
	err = postStorage.Create(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	posts, err := postStorage.ReadAllByBlog(context.Background(), blog.BlogID)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 {
		t.Fatal("expected one post linked to blog")
	}
}

func TestPostReadRecent(t *testing.T) {
	conn := connectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	postStorage := postgresql.NewPostStorage(conn)

	_, err := postStorage.ReadRecent(context.Background(), 20, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostReadSearch(t *testing.T) {
	conn := connectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	postStorage := postgresql.NewPostStorage(conn)

	_, err := postStorage.ReadSearch(context.Background(), "", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostCountRecent(t *testing.T) {
	conn := connectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	postStorage := postgresql.NewPostStorage(conn)

	_, err := postStorage.CountRecent(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostCountSearch(t *testing.T) {
	conn := connectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	postStorage := postgresql.NewPostStorage(conn)

	_, err := postStorage.CountSearch(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
}
