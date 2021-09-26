package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/task"
	"github.com/theandrew168/bloggulus/internal/web"
)

//go:embed migrations
var migrationsFS embed.FS

//go:embed static
var staticFS embed.FS

//go:embed templates
var templatesFS embed.FS

func main() {
	log.SetFlags(0)

	env := os.Getenv("ENV")
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

	migrations, _ := fs.Sub(migrationsFS, "migrations")
	if err = postgresql.Migrate(conn, context.Background(), migrations); err != nil {
		log.Fatal(err)
	}

	// reload templates from filesystem if ENV starts with "dev"
	var templates fs.FS
	if strings.HasPrefix(env, "dev") {
		templates = os.DirFS("templates")
	} else {
		templates, _ = fs.Sub(templatesFS, "templates")
	}

	static, _ := fs.Sub(staticFS, "static")

	app := &web.Application{
		StaticFS:    static,
		TemplatesFS: templates,

		Account: postgresql.NewAccountStorage(conn),
		Blog:    postgresql.NewBlogStorage(conn),
		Follow:  postgresql.NewFollowStorage(conn),
		Post:    postgresql.NewPostStorage(conn),
		Session: postgresql.NewSessionStorage(conn),
	}

	// kick off blog sync task
	syncBlogs := task.SyncBlogs(app.Blog, app.Post)
	go syncBlogs.Run(1 * time.Hour)

	// kick off session prune task
	pruneSessions := task.PruneSessions(app.Session)
	go pruneSessions.Run(5 * time.Minute)

	// kick off prometheus metrics handler
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("127.0.0.1:2112", mux)

	log.Printf("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router()))
}
