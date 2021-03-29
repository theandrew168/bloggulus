package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage/postgres"
	"github.com/theandrew168/bloggulus/tasks"
	"github.com/theandrew168/bloggulus/views"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mmcdole/gofeed"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:5000", "server listen address")
	addblog := flag.Bool("addblog", false, "-addblog <feed_url> <site_url>")
	syncblogs := flag.Bool("syncblogs", false, "sync blog posts with the database")
	flag.Parse()

	// ensure conn string env var exists
	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	// test a Connect and Ping now to verify DB connectivity
	db, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//	https://github.com/jackc/pgx/commit/aa8604b5c22989167e7158ecb1f6e7b8ddfebf04
	//	if err = db.Ping(ctx); err != nil {
	//		log.Fatal(err)
	//	}

	// apply database migrations
	if err = postgres.Migrate(db, context.Background(), "migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	// create storage interfaces
	accountStorage := postgres.NewAccountStorage(db)
	blogStorage := postgres.NewBlogStorage(db)
	postStorage := postgres.NewPostStorage(db)
	sessionStorage := postgres.NewSessionStorage(db)
	sourcedPostStorage := postgres.NewSourcedPostStorage(db)

	app := &views.Application{
		Account:     accountStorage,
		Blog:        blogStorage,
		Post:        postStorage,
		Session:     sessionStorage,
		SourcedPost: sourcedPostStorage,
	}

	if *addblog {
		feedURL := os.Args[2]
		siteURL := os.Args[3]
		log.Printf("adding blog: %s\n", feedURL)

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("  found: %s\n", feed.Title)

		blog := &models.Blog{
			FeedURL: feedURL,
			SiteURL: siteURL,
			Title:   feed.Title,
		}
		_, err = blogStorage.Create(context.Background(), blog)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *syncblogs {
		syncBlogs := tasks.SyncBlogs(blogStorage, postStorage)
		syncBlogs.RunNow()
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(app.HandleIndex))
	mux.Handle("/about", http.HandlerFunc(app.HandleAbout))
	mux.Handle("/blogs", http.HandlerFunc(app.HandleBlogs))
	mux.Handle("/login", http.HandlerFunc(app.HandleLogin))
	mux.Handle("/logout", http.HandlerFunc(app.HandleLogout))
	mux.Handle("/register", http.HandlerFunc(app.HandleRegister))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
