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

	"github.com/theandrew168/bloggulus/app"
	"github.com/theandrew168/bloggulus/rss"
	"github.com/theandrew168/bloggulus/storage"
	"github.com/theandrew168/bloggulus/storage/postgres"
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

		blog, err := rss.ReadBlog(feedURL)
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
		syncBlogs := task.SyncBlogs(blogStorage, postStorage)
		syncBlogs.RunNow()
		return
	}

	// kick off blog sync task
	syncBlogs := task.SyncBlogs(blogStorage, postStorage)
	go syncBlogs.Run(1 * time.Hour)

	// kick off session prune task
	pruneSessions := task.PruneSessions(sessionStorage)
	go pruneSessions.Run(5 * time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(app.HandleIndex))
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
