package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/app"
	"github.com/theandrew168/bloggulus/feeds"
	"github.com/theandrew168/bloggulus/storage"
	"github.com/theandrew168/bloggulus/storage/postgres"
	"github.com/theandrew168/bloggulus/tasks"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:5000", "server listen address")
	addblog := flag.Bool("addblog", false, "-addblog <feed_url>")
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

	if err = db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	// apply database migrations
	if err = postgres.Migrate(db, context.Background(), "migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	// create storage interfaces
	accountStorage := postgres.NewAccountStorage(db)
	blogStorage := postgres.NewBlogStorage(db)
	accountBlogStorage := postgres.NewAccountBlogStorage(db)
	postStorage := postgres.NewPostStorage(db)
	sessionStorage := postgres.NewSessionStorage(db)

	app := &app.Application{
		Account:     accountStorage,
		Blog:        blogStorage,
		AccountBlog: accountBlogStorage,
		Post:        postStorage,
		Session:     sessionStorage,
	}

	if *addblog {
		feedURL := os.Args[2]
		log.Printf("adding blog: %s\n", feedURL)

		blog, err := feeds.ReadBlog(feedURL)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("  found: %s\n", blog.Title)

		_, err = blogStorage.Create(context.Background(), blog)
		if err != nil {
			if err == storage.ErrDuplicateModel {
				log.Println("  already exists")
			} else {
				log.Fatal(err)
			}
		}
		return
	}

	if *syncblogs {
		syncBlogs := tasks.SyncBlogs(blogStorage, postStorage)
		syncBlogs.RunNow()
		return
	}

	// kick off blog sync task
	syncBlogs := tasks.SyncBlogs(blogStorage, postStorage)
	go syncBlogs.Run(1 * time.Hour)

	// kick off session prune task
	pruneSessions := tasks.PruneSessions(sessionStorage)
	go pruneSessions.Run(5 * time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(app.HandleIndex))
	mux.Handle("/about", http.HandlerFunc(app.HandleAbout))
	mux.Handle("/login", http.HandlerFunc(app.HandleLogin))
	mux.Handle("/logout", http.HandlerFunc(app.HandleLogout))
	mux.Handle("/blogs", http.HandlerFunc(app.HandleBlogs))
	mux.Handle("/follow", http.HandlerFunc(app.HandleFollow))
	mux.Handle("/unfollow", http.HandlerFunc(app.HandleUnfollow))
	mux.Handle("/register", http.HandlerFunc(app.HandleRegister))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
