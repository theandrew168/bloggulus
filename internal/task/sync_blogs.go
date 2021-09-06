package task

import (
	"context"
	"log"
	"sync"
	"time"

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
		go t.syncBlog(&wg, blog.BlogID, blog.FeedURL)
	}

	wg.Wait()
	return nil
}

func (t *syncBlogsTask) syncBlog(wg *sync.WaitGroup, blogID int, feedURL string) {
	defer wg.Done()

	// read current list of posts
	posts, err := feed.ReadPosts(feedURL)
	if err != nil {
		log.Println(err)
		return
	}

	// sync each post with the database
	for _, post := range posts {
		post.BlogID = blogID
		// TODO: get this from the feed or the page itself
		post.Preview = "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Tempora expedita dicta totam aspernatur doloremque. Excepturi iste iusto eos enim reprehenderit nisi, accusamus delectus nihil quis facere in modi ratione libero!"
		_, err := t.post.Create(context.Background(), post)
		if err != nil {
			if err != core.ErrExist {
				log.Println(err)
			}
		}
	}
}
