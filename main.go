package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"

	"github.com/theandrew168/bloggulus/app"
	"github.com/theandrew168/bloggulus/model"
	"github.com/theandrew168/bloggulus/model/postgres"
	"github.com/theandrew168/bloggulus/rss"
	"github.com/theandrew168/bloggulus/task"
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
	if err = migrate(db, context.Background(), "migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	// init app with storage interfaces
	app := &app.Application{
		Account:     postgres.NewAccountStorage(db),
		Blog:        postgres.NewBlogStorage(db),
		AccountBlog: postgres.NewAccountBlogStorage(db),
		Post:        postgres.NewPostStorage(db),
		Session:     postgres.NewSessionStorage(db),
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

	log.Printf("Listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, router))
}

func migrate(db *pgxpool.Pool, ctx context.Context, migrationsGlob string) error {
	// create migrations table if it doesn't exist
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			migration_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := db.Query(ctx, "SELECT name FROM migration")
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
	migrations, err := filepath.Glob(migrationsGlob)
	if err != nil {
		return err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		if _, ok := applied[migration]; !ok {
			missing = append(missing, migration)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)
	for _, file := range missing {
		log.Printf("applying: %s\n", file)

		// apply the missing ones
		sql, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		_, err = db.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = db.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", file)
		if err != nil {
			return err
		}
	}

	return nil
}
