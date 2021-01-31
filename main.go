package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/bmizerany/pat"
	"github.com/mmcdole/gofeed"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/acme/autocert"
)

type Blog struct {
	BlogID  int
	FeedURL string
	SiteURL string
	Title   string
}

type BlogStorage interface {
	Create(feedURL, siteURL, title string) (*Blog, error)
	Read(id int) (*Blog, error)
	ReadAll() ([]*Blog, error)
}

type Post struct {
	PostID  int
	BlogID  int
	URL     string
	Title   string
	Updated time.Time
}

type Application struct {
	ctx   context.Context
	db    *pgxpool.Pool
	index *template.Template
	about *template.Template
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			post.url,
			post.title,
			post.updated,
			blog.title
		FROM post
		JOIN blog ON blog.blog_id = post.blog_id
		ORDER BY post.updated DESC
		LIMIT 20`
	rows, err := app.db.Query(app.ctx, query)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type tmplPost struct {
		URL     string
		Title   string
		Updated time.Time
		Blog    string
	}

	var posts []*tmplPost
	for rows.Next() {
		var post tmplPost
		err = rows.Scan(&post.URL, &post.Title, &post.Updated, &post.Blog)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		posts = append(posts, &post)
	}

	data := struct{
		Posts  []*tmplPost
		Search string
	}{
		Posts:  posts,
		Search: "lol foobar",
	}

	ts, err := template.ParseFiles("templates/index.html", "templates/base.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) HandleAbout(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM blog ORDER BY title ASC"
	rows, err := app.db.Query(app.ctx, query)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var blog Blog
		err = rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		blogs = append(blogs, &blog)
	}

	data := struct {
		Blogs []*Blog
		Search string
	}{
		Blogs:  blogs,
		Search: "lol foobar",
	}

	ts, err := template.ParseFiles("templates/about.html", "templates/base.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) AddBlog(feedURL string, siteURL string) error {
	fp := gofeed.NewParser()
	blog, err := fp.ParseURL(feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("  found blog: %s\n", blog.Title)

	stmt := "INSERT INTO blog (feed_url, site_url, title) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	_, err = app.db.Exec(app.ctx, stmt, feedURL, siteURL, blog.Title)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) syncPost(wg *sync.WaitGroup, blogID int, post *gofeed.Item) {
	defer wg.Done()

	// use an old date if the post doesn't have one
	var updated time.Time
	if post.UpdatedParsed != nil {
		updated = *post.UpdatedParsed
	} else {
		updated = time.Now().AddDate(0, -1, 0)
	}

	stmt := "INSERT INTO post (blog_id, url, title, updated) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
	_, err := app.db.Exec(app.ctx, stmt, blogID, post.Link, post.Title, updated)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) syncBlog(wg *sync.WaitGroup, blogID int, url string) {
	defer wg.Done()

	fmt.Printf("checking blog: %s\n", url)

	// check if blog has been updated
	fp := gofeed.NewParser()
	blog, err := fp.ParseURL(url)
	if err != nil {
		log.Println(err)
		return
	}

	// sync each post in parallel
	for _, post := range blog.Items {
		fmt.Printf("updating post: %s\n", post.Title)
		wg.Add(1)
		go app.syncPost(wg, blogID, post)
	}
}

func (app *Application) SyncBlogs() error {
	// grab the current list of blogs
	query := "SELECT blog_id, feed_url FROM blog"
	rows, err := app.db.Query(app.ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// sync each blog in parallel
	var wg sync.WaitGroup
	for rows.Next() {
		var blogID int
		var url string
		err = rows.Scan(&blogID, &url)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("syncing blog: %s\n", url)

		wg.Add(1)
		go app.syncBlog(&wg, blogID, url)
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
	addblog := flag.Bool("addblog", false, "-addblog <feed_url> <site_url>")
	syncblogs := flag.Bool("syncblogs", false, "sync blog posts with the database")
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
	index := template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	about := template.Must(template.ParseFiles("templates/about.html", "templates/base.html"))

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

	if *addblog {
		fmt.Printf("adding blog: %s\n", os.Args[2])
		err = app.AddBlog(os.Args[2], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *syncblogs {
		fmt.Printf("syncing blogs\n")
		err = app.SyncBlogs()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.HandleIndex))
	mux.Get("/about", http.HandlerFunc(app.HandleAbout))

	// check if running via systemd
	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
		if os.Getenv("LISTEN_FDS") != "2" {
			log.Fatalf("expected 2 sockets from systemd, got %s\n", os.Getenv("LISTEN_FDS"))
		}

		// create http listener from the port 80 fd
		s80 := os.NewFile(3, "s80")
		httpListener, err := net.FileListener(s80)
		if err != nil {
			log.Fatal(err)
		}

		// create http listener from the port 443 fd
		s443 := os.NewFile(4, "s443")
		httpsListener, err := net.FileListener(s443)
		if err != nil {
			log.Fatal(err)
		}

		// redirect http to https
		go http.Serve(httpListener, http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}))

		// setup autocert to automatically obtain and renew TLS certs
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("bloggulus.com", "www.bloggulus.com"),
		}

		// create dir for autocert cache
		dir := filepath.Join(os.Getenv("HOME"), ".cache", "golang-autocert")
		if err := os.MkdirAll(dir, 0700); err != nil {
			log.Printf("warning: autocert not using cache: %v", err)
		} else {
			m.Cache = autocert.DirCache(dir)
		}

		s := &http.Server{
			Handler:      mux,
			TLSConfig:    m.TLSConfig(),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		fmt.Println("Listening on 0.0.0.0:80 and 0.0.0.0:443")
		s.ServeTLS(httpsListener, "", "")
	} else {
		fmt.Printf("Listening on %s\n", *addr)
		log.Fatal(http.ListenAndServe(*addr, mux))
	}
}
