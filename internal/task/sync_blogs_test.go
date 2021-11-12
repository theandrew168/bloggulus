package task_test

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/task"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestSyncBlogs(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)
	postStorage := postgresql.NewPostStorage(conn)

	// mock and create a blog
	blog := test.NewMockBlog()
	err := blogStorage.Create(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// mock some posts onto the blog
	posts := []core.Post{
		test.NewMockPost(blog),
		test.NewMockPost(blog),
		test.NewMockPost(blog),
	}
	body := test.RandomString(256)

	// create a feed reader for the mocked blog / post data
	reader := test.NewMockFeedReader(blog, posts, body)
	logger := log.New(io.Discard, "", 0)

	// run the sync blogs task
	syncBlogs := task.SyncBlogs(blogStorage, postStorage, reader, logger)
	err = syncBlogs.RunNow()
	if err != nil {
		t.Fatal(err)
	}

	// grab all posts associated with the mock blog
	synced, err := postStorage.ReadAllByBlog(context.Background(), blog.BlogID)
	if err != nil {
		t.Fatal(err)
	}

	// ensure that the posts were synced
	if len(synced) != len(posts) {
		t.Fatalf("want %v, got %v", len(posts), len(synced))
	}
}
