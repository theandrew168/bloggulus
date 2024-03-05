package task

import (
	"net/http"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain"
	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/storage"
)

type syncBlogsTask struct {
	w       *Worker
	storage *storage.Storage
	reader  feed.Reader
}

func (w *Worker) SyncBlogs(storage *storage.Storage, reader feed.Reader) Task {
	task := syncBlogsTask{
		w:       w,
		storage: storage,
		reader:  reader,
	}
	return &task
}

func (t *syncBlogsTask) Run(interval time.Duration) {
	err := t.RunNow()
	if err != nil {
		t.w.logger.Println(err)
	}

	c := time.Tick(interval)
	for {
		<-c

		err := t.syncBlogs()
		if err != nil {
			t.w.logger.Println(err)
		}
	}
}

func (t *syncBlogsTask) RunNow() error {
	return t.syncBlogs()
}

func (t *syncBlogsTask) syncBlogs() error {
	t.w.Add(1)
	defer t.w.Done()

	limit := 50
	offset := 0

	// read initial batch of blogs
	blogs, err := t.storage.Blog.ReadAll(limit, offset)
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
		blogs, err = t.storage.Blog.ReadAll(limit, offset)
		if err != nil {
			wg.Wait()
			return err
		}
	}

	wg.Wait()
	return nil
}

func (t *syncBlogsTask) syncBlog(wg *sync.WaitGroup, blog domain.Blog) {
	defer wg.Done()

	req, err := http.NewRequest("GET", blog.FeedURL, nil)
	if err != nil {
		t.w.logger.Printf("%d: %s\n", blog.ID, err)
		return
	}

	if blog.ETag != "" {
		req.Header.Set("If-None-Match", blog.ETag)
	}
	if blog.LastModified != "" {
		req.Header.Set("If-Modified-Since", blog.LastModified)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.w.logger.Printf("%d: %s\n", blog.ID, err)
		return
	}
	defer resp.Body.Close()

	etag := resp.Header.Get("ETag")
	lastModified := resp.Header.Get("Last-Modified")

	dirty := false
	if etag != "" && etag != blog.ETag {
		blog.ETag = etag
		dirty = true
	}
	if lastModified != "" && lastModified != blog.LastModified {
		blog.LastModified = lastModified
		dirty = true
	}

	// update blog if ETag or Last-Modified have changed
	if dirty {
		err = t.storage.Blog.Update(blog)
		if err != nil {
			t.w.logger.Printf("%d: %s\n", blog.ID, err)
			return
		}
	}

	// early exit for "no new content" or errors
	if resp.StatusCode >= 300 {
		t.w.logger.Printf("%d: skipping: code: %d\n", blog.ID, resp.StatusCode)
		return
	}

	limit := 50
	offset := 0

	// build a set of known post URLs
	knownPostURLs := make(map[string]bool)

	// read initial batch of posts
	knownPosts, err := t.storage.Post.ReadAllByBlog(blog, limit, offset)
	if err != nil {
		t.w.logger.Printf("%d: %s\n", blog.ID, err)
		return
	}

	for len(knownPosts) > 0 {
		// add each post URL to the set
		for _, post := range knownPosts {
			knownPostURLs[post.URL] = true
		}

		// read the next batch
		offset += limit
		knownPosts, err = t.storage.Post.ReadAllByBlog(blog, limit, offset)
		if err != nil {
			t.w.logger.Printf("%d: %s\n", blog.ID, err)
			return
		}
	}

	// read posts from feed
	feedPosts, err := t.reader.ReadBlogPosts(blog, resp.Body)
	if err != nil {
		t.w.logger.Printf("%d: %s\n", blog.ID, err)
		return
	}

	// newPosts = feedPosts - knownPosts
	var newPosts []domain.Post
	for _, post := range feedPosts {
		if _, ok := knownPostURLs[post.URL]; ok {
			continue
		}
		newPosts = append(newPosts, post)
	}

	// attempt to read the content for each new post
	for i := range newPosts {
		t.w.logger.Printf("%d: sync: %v\n", blog.ID, newPosts[i].URL)
		body, err := t.reader.ReadPostBody(newPosts[i])
		if err != nil {
			t.w.logger.Printf("%d: %s\n", blog.ID, err)
			continue
		}
		newPosts[i].Body = body
	}

	// sync each post with the database
	for _, post := range newPosts {
		err = t.storage.Post.Create(&post)
		if err != nil {
			t.w.logger.Printf("%d: sync: %v %v\n", blog.ID, post.URL, err)
		}
	}
}
