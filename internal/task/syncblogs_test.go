package task_test

import (
	"testing"

	"github.com/theandrew168/bloggulus"
	"github.com/theandrew168/bloggulus/internal/feed"
	"github.com/theandrew168/bloggulus/internal/task"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestSyncBlogs(t *testing.T) {
	logger := test.NewLogger(t)
	storage, closer := test.NewStorage(t)
	defer closer()

	// mock and create a blog
	blog := test.NewMockBlog()
	err := storage.Blog.Create(&blog)
	if err != nil {
		t.Fatal(err)
	}

	// mock some posts onto the blog
	posts := []bloggulus.Post{
		test.NewMockPost(blog),
		test.NewMockPost(blog),
		test.NewMockPost(blog),
	}
	body := test.RandomString(256)

	// create a feed reader for the mocked blog / post data
	reader := feed.NewMockReader(blog, posts, body)

	// run the sync blogs task
	worker := task.NewWorker(logger)
	syncBlogs := worker.SyncBlogs(storage, reader)
	err = syncBlogs.RunNow()
	if err != nil {
		t.Fatal(err)
	}

	// grab all posts associated with the mock blog
	synced, err := storage.Post.ReadAllByBlog(blog, 20, 0)
	if err != nil {
		t.Fatal(err)
	}

	// ensure that the posts were synced
	if len(synced) != len(posts) {
		t.Fatalf("want %v, got %v\n", len(posts), len(synced))
	}
}
