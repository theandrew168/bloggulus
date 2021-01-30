package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/bmizerany/pat"
	"github.com/mmcdole/gofeed"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Application struct {
	ctx   context.Context
	db    *pgxpool.Pool
	index *template.Template
	about *template.Template
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) HandleAbout(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) AddFeed(url string, siteUrl string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	fmt.Printf("  found feed: %s\n", feed.Title)

	stmt := "INSERT INTO feed (url, site_url, title) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	_, err = app.db.Exec(app.ctx, stmt, url, siteUrl, feed.Title)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) syncPost(wg *sync.WaitGroup, feedId int, post *gofeed.Item) {
	defer wg.Done()

	// use an old date if the post doesn't have one
	var updated time.Time
	if post.UpdatedParsed != nil {
		updated = *post.UpdatedParsed
	} else {
		updated = time.Now().AddDate(0, -1, 0)
	}

	stmt := "INSERT INTO post (feed_id, url, title, updated) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	_, err := app.db.Exec(app.ctx, stmt, feedId, post.Link, post.Title, updated)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) syncFeed(wg *sync.WaitGroup, feedId int, url string) {
	defer wg.Done()

	fmt.Printf("checking feed: %s\n", url)

	// check if feed has been updated
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Println(err)
		return
	}

	// sync each post in parallel
	for _, post := range feed.Items {
		fmt.Printf("updating post: %s\n", post.Title)
		wg.Add(1)
		go app.syncPost(wg, feedId, post)
	}
}

func (app *Application) SyncFeeds() error {
	// grab the current list of feeds
	query := "SELECT feed_id, url FROM feed"
	rows, err := app.db.Query(app.ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// sync each feed in parallel
	var wg sync.WaitGroup
	for rows.Next() {
		var feedId int
		var url string
		err = rows.Scan(&feedId, &url)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("syncing feed: %s\n", url)

		wg.Add(1)
		go app.syncFeed(&wg, feedId, url)
	}

	wg.Wait()
	return nil
}

func (app *Application) Migrate(migrationsGlob string) error {
	// create migration table if it doesn't exist
	_, err := app.db.Exec(app.ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			migration_id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := app.db.Query(app.ctx, "SELECT name FROM migration")
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
		fmt.Printf("applying: %s\n", file)

		// apply the missing ones
		sql, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		_, err = app.db.Exec(app.ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = app.db.Exec(app.ctx, "INSERT INTO migration (name) VALUES ($1)", file)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "server listen address")
	addfeed := flag.Bool("addfeed", false, "-addfeed <url> <site_url>")
	syncfeeds := flag.Bool("syncfeeds", false, "sync feeds with the database")
	flag.Parse()

	// ensure conn string env var exists
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Missing required env var: DATABASE_URL")
	}

	// test a Connect and Ping now to verify DB connectivity
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

//	if err = db.Ping(ctx); err != nil {
//		log.Fatal(err)
//	}

	// load the necessary templates
	index := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	about := template.Must(template.ParseFiles("templates/base.html", "templates/about.html"))

	app := &Application{
		ctx:   ctx,
		db:    db,
		index: index,
		about: about,
	}

	// apply database migrations
	if err = app.Migrate("migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	if *addfeed {
		fmt.Printf("adding feed: %s\n", os.Args[2])
		err = app.AddFeed(os.Args[2], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *syncfeeds {
		fmt.Printf("syncing feeds\n")
		err = app.SyncFeeds()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.HandleIndex))
	mux.Get("/about", http.HandlerFunc(app.HandleAbout))

	fmt.Printf("Listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
