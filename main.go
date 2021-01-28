package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/bmizerany/pat"
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
	db, err := sql.Open("sqlite3", "bloggulus.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	err = migrate(db, "migrations/*.sql")
	if err != nil {
		log.Fatal(err)
	}

	index := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	about := template.Must(template.ParseFiles("templates/base.html", "templates/about.html"))

	app := &Application{
		db:    db,
		index: index,
		about: about,
	}

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.HandleIndex))
	mux.Get("/about", http.HandlerFunc(app.HandleAbout))

	addr := "localhost:8080"
	fmt.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
