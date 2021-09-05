package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"

	"github.com/theandrew168/bloggulus/internal/app"
	"github.com/theandrew168/bloggulus/internal/model"
	"github.com/theandrew168/bloggulus/internal/model/postgresql"
	"github.com/theandrew168/bloggulus/internal/rss"
	"github.com/theandrew168/bloggulus/internal/task"
)

func main() {
	addblog := flag.Bool("addblog", false, "-addblog <feed_url>")
	syncblogs := flag.Bool("syncblogs", false, "sync blog posts with the database")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	db, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	// init app with storage interfaces
	app := &app.Application{
		Account:     postgresql.NewAccountStorage(db),
		Blog:        postgresql.NewBlogStorage(db),
		AccountBlog: postgresql.NewAccountBlogStorage(db),
		Post:        postgresql.NewPostStorage(db),
		Session:     postgresql.NewSessionStorage(db),
	}

	if *addblog {
		feedURL := os.Args[2]
		log.Printf("adding blog: %s\n", feedURL)

		blog, err := rss.ReadBlog(feedURL)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("  found: %s\n", blog.Title)

		_, err = app.Blog.Create(context.Background(), blog)
		if err != nil {
			if err == model.ErrExist {
				log.Println("  already exists")
			} else {
				log.Fatal(err)
			}
		}
		return
	}

	if *syncblogs {
		syncBlogs := task.SyncBlogs(app.Blog, app.Post)
		syncBlogs.RunNow()
		return
	}

	// kick off blog sync task
	syncBlogs := task.SyncBlogs(app.Blog, app.Post)
	go syncBlogs.Run(1 * time.Hour)

	// kick off session prune task
	pruneSessions := task.PruneSessions(app.Session)
	go pruneSessions.Run(5 * time.Minute)

	router := httprouter.New()
	router.HandlerFunc("GET", "/", app.HandleIndex)
	router.HandlerFunc("GET", "/login", app.HandleLogin)
	router.HandlerFunc("POST", "/login", app.HandleLogin)
	router.HandlerFunc("POST", "/logout", app.HandleLogout)
	router.HandlerFunc("GET", "/blogs", app.HandleBlogs)
	router.HandlerFunc("POST", "/blogs", app.HandleBlogs)
	router.HandlerFunc("POST", "/follow", app.HandleFollow)
	router.HandlerFunc("POST", "/unfollow", app.HandleUnfollow)
	router.HandlerFunc("GET", "/register", app.HandleRegister)
	router.HandlerFunc("POST", "/register", app.HandleRegister)
	router.ServeFiles("/static/*filepath", http.Dir("./static"))

	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
