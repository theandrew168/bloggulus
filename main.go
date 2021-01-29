package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/mmcdole/gofeed"
	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	db    *sql.DB
	index *template.Template
	about *template.Template
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) HandleAbout(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) AddFeed(url, siteUrl string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	fmt.Printf("  found feed: %s\n", feed.Title)

	stmt := "INSERT INTO feed (url, site_url, title) VALUES (?, ?, ?)"
	_, err = app.db.Exec(stmt, url, siteUrl, feed.Title)
	if err != nil {
		// move along if the feed is already here
		if strings.Contains(err.Error(), "UNIQUE") {
			fmt.Println("  feed already exists")
			return nil
		}
		// else the error is something worth looking at
		return err
	}
	return nil
}

func migrate(db *sql.DB, migrationsGlob string) error {
	// create migration table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migration (
			migration_id INTEGER PRIMARY KEY NOT NULL,
			name TEXT UNIQUE NOT NULL
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := db.Query("SELECT name FROM migration")
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
		_, err = db.Exec(string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = db.Exec("INSERT INTO migration (name) VALUES (?)", file)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	addr := flag.String("addr", "localhost:8080", "listen address")
	addfeed := flag.Bool("addfeed", false, "-addfeed <url> <site_url>")
	flag.Parse()

	db, err := sql.Open("sqlite3", "bloggulus.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err = migrate(db, "migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	index := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	about := template.Must(template.ParseFiles("templates/base.html", "templates/about.html"))

	app := &Application{
		db:    db,
		index: index,
		about: about,
	}

	if *addfeed {
		fmt.Printf("adding feed: %s\n", os.Args[2])
		err = app.AddFeed(os.Args[2], os.Args[3])
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
