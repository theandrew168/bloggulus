package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/model/postgresql"
	"github.com/theandrew168/bloggulus/internal/task"
	"github.com/theandrew168/bloggulus/internal/web"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	// init app with storage interfaces
	app := &web.Application{
		Account:     postgresql.NewAccountStorage(conn),
		Blog:        postgresql.NewBlogStorage(conn),
		AccountBlog: postgresql.NewAccountBlogStorage(conn),
		Post:        postgresql.NewPostStorage(conn),
		Session:     postgresql.NewSessionStorage(conn),
	}

	// kick off blog sync task
	syncBlogs := task.SyncBlogs(app.Blog, app.Post)
	go syncBlogs.Run(1 * time.Hour)

	// kick off session prune task
	pruneSessions := task.PruneSessions(app.Session)
	go pruneSessions.Run(5 * time.Minute)

	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router()))
}
