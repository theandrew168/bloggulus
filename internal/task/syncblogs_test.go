package task_test

import (
	"context"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/feed"
	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/task"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestSyncBlogs(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	// instantiate storage interface
	storage := postgresql.NewStorage(conn)

	// mock and create a blog
	blog := test.NewMockBlog()
	err := storage.BlogCreate(context.Background(), &blog)
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
	reader := feed.NewMockReader(blog, posts, body)
	logger := test.NewLogger()

	// run the sync blogs task
	syncBlogs := task.SyncBlogs(storage, reader, logger)
	err = syncBlogs.RunNow()
	if err != nil {
		t.Fatal(err)
	}

	// grab all posts associated with the mock blog
	synced, err := storage.PostReadAllByBlog(context.Background(), blog.ID)
	if err != nil {
		t.Fatal(err)
	}

	// ensure that the posts were synced
	if len(synced) != len(posts) {
		t.Fatalf("want %v, got %v\n", len(posts), len(synced))
	}
}
