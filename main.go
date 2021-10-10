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

	"github.com/jackc/pgx/v4/pgxpool"

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
	log.SetFlags(0)

	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	addblog := flag.Bool("addblog", false, "-addblog <feed_url>")
	flag.Parse()

	databaseURL := os.Getenv("BLOGGULUS_DATABASE_URL")
	if databaseURL == "" {
		log.Fatalln("Missing required env var: BLOGGULUS_DATABASE_URL")
	}

	conn, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatalln(err)
	}

	migrations, _ := fs.Sub(migrationsFS, "migrations")
	if err = postgresql.Migrate(conn, context.Background(), migrations); err != nil {
		log.Fatalln(err)
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

		Blog: postgresql.NewBlogStorage(conn),
		Post: postgresql.NewPostStorage(conn),
	}

//	post, err := app.Post.Read(context.Background(), 1)
//	if err == nil {
//		log.Printf("%v %v\n", post.Title, post.Tags)
//	}

	if *addblog {
		feedURL := os.Args[2]
		log.Printf("adding blog: %s\n", feedURL)

		blog, err := feed.ReadBlog(feedURL)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("  found: %s\n", blog.Title)

		err = app.Blog.Create(context.Background(), &blog)
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
	syncBlogs := task.SyncBlogs(app.Blog, app.Post)
	go syncBlogs.Run(1 * time.Hour)

	log.Printf("listening on %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, app.Router()))
}
