package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/bloggulus/internal/api"
	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/feed"
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
	// silence timestamp and log level
	log.SetFlags(0)

	// check for general config vars
	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	// check for flags
	addblog := flag.Bool("addblog", false, "-addblog <feed_url>")
	flag.Parse()

	// check for database connection url var
	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		log.Fatalln("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// apply database migrations
	migrations, _ := fs.Sub(migrationsFS, "migrations")
	if err = postgresql.Migrate(conn, context.Background(), migrations); err != nil {
		log.Fatalln(err)
	}

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)
	postStorage := postgresql.NewPostStorage(conn)

	// add a blog and exit now if requested
	if *addblog {
		feedURL := os.Args[2]
		log.Printf("adding blog: %s\n", feedURL)

		blog, err := feed.ReadBlog(feedURL)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("  found: %s\n", blog.Title)

		err = blogStorage.Create(context.Background(), &blog)
		if err != nil {
			if err == core.ErrExist {
				log.Println("  already exists")
			} else {
				log.Fatal(err)
			}
		}

		return
	}

	// kick off blog sync task
	syncBlogs := task.SyncBlogs(blogStorage, postStorage)
	go syncBlogs.Run(1 * time.Hour)

	// reload templates from filesystem if ENV starts with "dev"
	var templates fs.FS
	if strings.HasPrefix(env, "dev") {
		templates = os.DirFS("templates")
	} else {
		templates, _ = fs.Sub(templatesFS, "templates")
	}

	// init web application struct
	webApp := &web.Application{
		TemplatesFS: templates,

		Blog: blogStorage,
		Post: postStorage,
	}

	// init api application struct
	apiApp := &api.Application{
		Blog: blogStorage,
		Post: postStorage,
	}

	// setup http.Handler for static files
	static, _ := fs.Sub(staticFS, "static")
	staticServer := http.FileServer(http.FS(static))

	// construct the top-level router
	r := chi.NewRouter()
	r.Mount("/", webApp.Router())
	r.Mount("/api", apiApp.Router())
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/static/*", http.StripPrefix("/static", staticServer))

	// lets go!
	log.Printf("listening on %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
