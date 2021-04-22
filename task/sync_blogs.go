package task

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/rss"
	"github.com/theandrew168/bloggulus/storage"
)

type syncBlogsTask struct {
	Blog *storage.Blog
	Post *storage.Post
}

func SyncBlogs(blog *storage.Blog, post *storage.Post) Task {
	return &syncBlogsTask{
		Blog: blog,
		Post: post,
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
	blogs, err := t.Blog.ReadAll(context.Background())
	if err != nil {
		return err
	}

	// sync each blog in parallel
	var wg sync.WaitGroup
	for _, blog := range blogs {
		wg.Add(1)
		go t.syncBlog(&wg, blog.BlogID, blog.FeedURL)
	}

	wg.Wait()
	return nil
}

func (t *syncBlogsTask) syncBlog(wg *sync.WaitGroup, blogID int, feedURL string) {
	defer wg.Done()

	log.Printf("syncing blog: %s\n", feedURL)

	// read current list of posts
	posts, err := rss.ReadPosts(feedURL)
	if err != nil {
		log.Println(err)
		return
	}

	// sync each post with the database
	for _, post := range posts {
		log.Printf("updating post: %s\n", post.Title)

		post.BlogID = blogID
		_, err := t.Post.Create(context.Background(), post)
		if err != nil {
			if err == storage.ErrDuplicateModel {
				log.Println("  already exists")
			} else {
				log.Println(err)
			}
		}
	}
}
