package task

import (
	"context"
	"html"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/feed"
)

type syncBlogsTask struct {
	blog core.BlogStorage
	post core.PostStorage
}

func SyncBlogs(blog core.BlogStorage, post core.PostStorage) Task {
	return &syncBlogsTask{
		blog: blog,
		post: post,
	}
}

func (t *syncBlogsTask) Run(interval time.Duration) {
	c := time.Tick(interval)
	for {
		<-c

		err := t.syncBlogs()
		if err != nil {
			log.Println(err)
		}
	}
}

func (t *syncBlogsTask) RunNow() error {
	return t.syncBlogs()
}

func (t *syncBlogsTask) syncBlogs() error {
	blogs, err := t.blog.ReadAll(context.Background())
	if err != nil {
		return err
	}

	// sync each blog in parallel
	var wg sync.WaitGroup
	for _, blog := range blogs {
		wg.Add(1)
		go t.syncBlog(&wg, blog)
	}

	wg.Wait()
	return nil
}

func (t *syncBlogsTask) syncBlog(wg *sync.WaitGroup, blog core.Blog) {
	defer wg.Done()

	// read posts currently in storage
	knownPosts, err := t.post.ReadAllByBlog(context.Background(), blog.BlogID)
	if err != nil {
		log.Println(err)
		return
	}

	// build a set of known post URLs
	knownPostURLs := make(map[string]bool)
	for _, post := range knownPosts {
		knownPostURLs[post.URL] = true
	}

	// read posts from feed
	feedPosts, err := feed.ReadPosts(blog)
	if err != nil {
		log.Println(err)
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

	// sync each post with the database
	for _, post := range newPosts {
		// attempt to manually fetch post body if not already present in feed
		if post.Body == "" {
			resp, err := http.Get(post.URL)
			if err != nil {
				log.Println(err)
				goto create
			}
			defer resp.Body.Close()

			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				goto create
			}

			p := bluemonday.StripTagsPolicy()
			body := string(buf)
			body = p.Sanitize(body)
			body = html.UnescapeString(body)

			post.Body = body
		}

create:
		err := t.post.Create(context.Background(), &post)
		if err != nil {
			if err != core.ErrExist {
				log.Println(err)
			}
		}
	}
}
