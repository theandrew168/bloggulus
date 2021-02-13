package tasks

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"

	"github.com/mmcdole/gofeed"
)

type syncBlogsTask struct {
	Blog storage.Blog
	Post storage.Post
}

func SyncBlogs(blog storage.Blog, post storage.Post) Task {
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
		log.Printf("syncing blog: %s\n", blog.FeedURL)

		wg.Add(1)
		go t.syncBlog(&wg, blog.BlogID, blog.FeedURL)
	}

	wg.Wait()
	return nil
}

func (t *syncBlogsTask) syncBlog(wg *sync.WaitGroup, blogID int, url string) {
	defer wg.Done()

	log.Printf("checking blog: %s\n", url)

	// check if blog has been updated
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Println(err)
		return
	}

	// sync each post in parallel
	for _, item := range feed.Items {
		log.Printf("updating post: %s\n", item.Title)
		wg.Add(1)
		go t.syncPost(wg, blogID, item)
	}
}

func (t *syncBlogsTask) syncPost(wg *sync.WaitGroup, blogID int, item *gofeed.Item) {
	defer wg.Done()

	// use an old date if the post doesn't have one
	var updated time.Time
	if item.UpdatedParsed != nil {
		updated = *item.UpdatedParsed
	} else {
		updated = time.Now().AddDate(0, -3, 0)
	}

	post := &models.Post{
		BlogID:  blogID,
		URL:     item.Link,
		Title:   item.Title,
		Updated: updated,
	}
	_, err := t.Post.Create(context.Background(), post)
	if err != nil {
		log.Println(err)
		return
	}
}
