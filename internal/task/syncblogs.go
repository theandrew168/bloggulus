package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/feed"
)

type syncBlogsTask struct {
	worker  *Worker
	storage core.Storage
	reader  feed.Reader
}

func (w *Worker) SyncBlogs(storage core.Storage, reader feed.Reader) Task {
	task := syncBlogsTask{
		worker:  w,
		storage: storage,
		reader:  reader,
	}
	return &task
}

func (t *syncBlogsTask) Run(interval time.Duration) {
	err := t.RunNow()
	if err != nil {
		t.worker.logError(err)
	}

	c := time.Tick(interval)
	for {
		<-c

		err := t.syncBlogs()
		if err != nil {
			t.worker.logError(err)
		}
	}
}

func (t *syncBlogsTask) RunNow() error {
	return t.syncBlogs()
}

func (t *syncBlogsTask) syncBlogs() error {
	t.worker.Add(1)
	defer t.worker.Done()

	limit := 50
	offset := 0

	// read initial batch of blogs
	blogs, err := t.storage.ReadBlogs(context.Background(), limit, offset)
	if err != nil {
		return err
	}

	// kick off blog syncs in batches
	var wg sync.WaitGroup
	for len(blogs) > 0 {
		// sync each blog in parallel
		for _, blog := range blogs {
			wg.Add(1)
			go t.syncBlog(&wg, blog)
		}

		// read the next batch
		offset += limit
		blogs, err = t.storage.ReadBlogs(context.Background(), limit, offset)
		if err != nil {
			wg.Wait()
			return err
		}
	}

	wg.Wait()
	return nil
}

func (t *syncBlogsTask) syncBlog(wg *sync.WaitGroup, blog core.Blog) {
	defer wg.Done()

	limit := 50
	offset := 0

	// build a set of known post URLs
	knownPostURLs := make(map[string]bool)

	// read initial batch of posts
	knownPosts, err := t.storage.ReadPostsByBlog(context.Background(), blog.ID, limit, offset)
	if err != nil {
		t.worker.logError(err)
		return
	}

	for len(knownPosts) > 0 {
		// add each post URL to the set
		for _, post := range knownPosts {
			knownPostURLs[post.URL] = true
		}

		// read the next batch
		offset += limit
		knownPosts, err = t.storage.ReadPostsByBlog(context.Background(), blog.ID, limit, offset)
		if err != nil {
			t.worker.logError(err)
			return
		}
	}

	// read posts from feed
	feedPosts, err := t.reader.ReadBlogPosts(blog)
	if err != nil {
		t.worker.logError(err)
		return
	}

	// newPosts = feedPosts - knownPosts
	var newPosts []core.Post
	for _, post := range feedPosts {
		if _, ok := knownPostURLs[post.URL]; ok {
			continue
		}
		newPosts = append(newPosts, post)
	}

	// attempt to read the content for each new post
	for i, _ := range newPosts {
		body, err := t.reader.ReadPostBody(newPosts[i])
		if err != nil {
			t.worker.logError(err)
			continue
		}
		newPosts[i].Body = body
	}

	// sync each post with the database
	for _, post := range newPosts {
		err = t.storage.CreatePost(context.Background(), &post)
		if err != nil {
			msg := fmt.Sprintf("sync %v %v\n", post.URL, err)
			t.worker.log(msg)
		}
	}
}
