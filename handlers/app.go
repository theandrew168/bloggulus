package handlers

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage"

	"github.com/mmcdole/gofeed"
)

type Application struct {
	Account     storage.Account
	Blog        storage.Blog
	Post        storage.Post
	Session     storage.Session
	SourcedPost storage.SourcedPost
}

func (app *Application) syncPost(wg *sync.WaitGroup, blogID int, item *gofeed.Item) {
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
	_, err := app.Post.Create(context.Background(), post)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) syncBlog(wg *sync.WaitGroup, blogID int, url string) {
	defer wg.Done()

	fmt.Printf("checking blog: %s\n", url)

	// check if blog has been updated
	fp := gofeed.NewParser()
	blog, err := fp.ParseURL(url)
	if err != nil {
		log.Println(err)
		return
	}

	// sync each post in parallel
	for _, post := range blog.Items {
		fmt.Printf("updating post: %s\n", post.Title)
		wg.Add(1)
		go app.syncPost(wg, blogID, post)
	}
}

func (app *Application) SyncBlogs() error {
	blogs, err := app.Blog.ReadAll(context.Background())
	if err != nil {
		return err
	}

	// sync each blog in parallel
	var wg sync.WaitGroup
	for _, blog := range blogs {
		fmt.Printf("syncing blog: %s\n", blog.FeedURL)

		wg.Add(1)
		go app.syncBlog(&wg, blog.BlogID, blog.FeedURL)
	}

	wg.Wait()
	return nil
}

func (app *Application) HourlySync() {
	c := time.Tick(1 * time.Hour)
	for {
		<-c

		err := app.SyncBlogs()
		if err != nil {
			log.Println(err)
		}
	}
}
