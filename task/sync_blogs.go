package task

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/model"
	"github.com/theandrew168/bloggulus/rss"
)

type syncBlogsTask struct {
	Blog model.BlogStorage
	Post model.PostStorage
}

func SyncBlogs(blog model.BlogStorage, post model.PostStorage) Task {
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

	// read current list of posts
	posts, err := rss.ReadPosts(feedURL)
	if err != nil {
		log.Println(err)
		return
	}

	// sync each post with the database
	for _, post := range posts {
		post.BlogID = blogID
		_, err := t.Post.Create(context.Background(), post)
		if err != nil {
			if err != model.ErrExist {
				log.Println(err)
			}
		}
	}
}