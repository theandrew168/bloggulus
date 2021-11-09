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
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/klauspost/compress/gzhttp"
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
	logger := log.New(os.Stdout, "", 0)

	// check for general config vars
	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	// check for flags
	migrate := flag.Bool("migrate", false, "-migrate")
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
	if err = applyMigrations(conn, migrations, logger); err != nil {
		log.Fatalln(err)
	}

	// exit now if just applying migrations
	if *migrate {
		return
	}

	// init storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)
	postStorage := postgresql.NewPostStorage(conn)

	// init default feed reader
	reader := feed.NewReader()

	// add a blog and exit now if requested
	if *addblog {
		feedURL := os.Args[2]
		log.Printf("adding blog: %s\n", feedURL)

		blog, err := reader.ReadBlog(feedURL)
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
	syncBlogs := task.SyncBlogs(blogStorage, postStorage, reader, logger)
	go syncBlogs.Run(1 * time.Hour)

	// reload templates from filesystem if ENV starts with "dev"
	var templates fs.FS
	if strings.HasPrefix(env, "dev") {
		templates = os.DirFS("templates")
	} else {
		templates, _ = fs.Sub(templatesFS, "templates")
	}

	// init web application
	webApp := web.NewApplication(
		templates,
		blogStorage,
		postStorage,
		logger,
	)

	// init api application struct
	apiApp := api.NewApplication(
		blogStorage,
		postStorage,
		logger,
	)

	// setup http.Handler for static files
	static, _ := fs.Sub(staticFS, "static")
	staticServer := http.FileServer(http.FS(static))
	gzipStaticServer := gzhttp.GzipHandler(staticServer)

	// construct the top-level router
	r := chi.NewRouter()
	r.Mount("/", webApp.Router())
	r.Mount("/api", apiApp.Router())
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/static/*", http.StripPrefix("/static", gzipStaticServer))

	// lets go!
	log.Printf("listening on %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}

func applyMigrations(conn *pgxpool.Pool, migrationsFS fs.FS, logger *log.Logger) error {
	ctx := context.Background()

	// create migrations table if it doesn't exist
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			migration_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := conn.Query(ctx, "SELECT name FROM migration")
	if err != nil {
		return err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		applied[name] = true
	}

	// get migrations that should be applied (from migrations/ dir)
	migrations, err := fs.ReadDir(migrationsFS, ".")
	if err != nil {
		return err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		if _, ok := applied[name]; !ok {
			missing = append(missing, name)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)
	for _, name := range missing {
		logger.Printf("applying: %s\n", name)

		// apply the missing ones
		sql, err := fs.ReadFile(migrationsFS, name)
		if err != nil {
			return err
		}
		_, err = conn.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = conn.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}
	}

	logger.Printf("migrations up to date\n")
	return nil
}
