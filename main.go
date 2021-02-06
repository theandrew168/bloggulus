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

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/mmcdole/gofeed"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/acme/autocert"
)

// models
type Blog struct {
	BlogID  int
	FeedURL string
	SiteURL string
	Title   string
}

type Post struct {
	PostID  int
	BlogID  int
	URL     string
	Title   string
	Updated time.Time
}

// storage
type BloggulusStorage interface {
	CreateBlog(feedURL, siteURL, title string) (*Blog, error)
	ReadAllBlogs() ([]*Blog, error)

	CreatePost(blogID int, URL, title string, updated time.Time) (*Post, error)
	ReadRecentPosts(n int) ([]*Post, error)
}

// storage - postgres
type postgresStorage struct {
	ctx   context.Context
	db    *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, db *pgxpool.Pool) BloggulusStorage {
	return &postgresStorage{
		ctx: ctx,
		db:  db,
	}
}

// storage - postgres - blog
func (s *postgresStorage) CreateBlog(feedURL, siteURL, title string) (*Blog, error) {
	stmt := "INSERT INTO blog (feed_url, site_url, title) VALUES ($1, $2, $3) RETURNING blog_id"
	row := s.db.QueryRow(s.ctx, stmt, feedURL, siteURL, title)

	var blogID int
	err := row.Scan(&blogID)
	if err != nil {
		return nil, err
	}

	blog := &Blog{
		BlogID:  blogID,
		FeedURL: feedURL,
		SiteURL: siteURL,
		Title:   title,
	}

	return blog, nil
}

func (s *postgresStorage) ReadAllBlogs() ([]*Blog, error) {
	query := "SELECT * FROM blog"
	rows, err := s.db.Query(s.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.BlogID, &blog.FeedURL, &blog.SiteURL, &blog.Title)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, &blog)
	}

	return blogs, nil
}

// storage - postgres - post
func (s *postgresStorage) CreatePost(blogID int, URL, title string, updated time.Time) (*Post, error) {
	stmt := "INSERT INTO post (blog_id, url, title, updated) VALUES ($1, $2, $3, $4) RETURNING post_id"
	row := s.db.QueryRow(s.ctx, stmt, blogID, URL, title, updated)

	var postID int
	err := row.Scan(&postID)
	if err != nil {
		return nil, err
	}

	post := &Post{
		PostID:  postID,
		BlogID:  blogID,
		URL:     URL,
		Title:   title,
		Updated: updated,
	}

	return post, nil
}

func (s *postgresStorage) ReadRecentPosts(n int) ([]*Post, error) {
	query := "SELECT * FROM post ORDER BY updated DESC LIMIT $1"
	rows, err := s.db.Query(s.ctx, query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.BlogID, &post.URL, &post.Title, &post.Updated)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

// app stuff
type Application struct {
	session *scs.SessionManager
	store   BloggulusStorage
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) HandleAbout(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("attempted login from user: %s\n", r.PostFormValue("username"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("attempted register from user: %s\n", r.PostFormValue("username"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
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

	_, err := app.store.CreatePost(blogID, post.Link, post.Title, updated)
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
	blogs, err := app.store.ReadAllBlogs()
	if err != nil {
		return err
	}

	// sync each blog in parallel
	var wg sync.WaitGroup
	for _, blog := range blogs {
		fmt.Printf("syncing blog: %s\n", blog.FeedURL)

		wg.Add(1)
		go app.syncBlog(&wg, blog.BlogID, blog.FeedURL)
	}

	wg.Wait()
	return nil
}

func (app *Application) HourlySync() {
	c := time.Tick(1 * time.Hour)
	for {
		<-c

		err := app.SyncBlogs()
		if err != nil {
			log.Println(err)
		}
	}
}

func migrate(ctx context.Context, db *pgxpool.Pool, migrationsGlob string) error {
	// create migration table if it doesn't exist
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			migration_id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
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
		fmt.Printf("applying: %s\n", file)

		// apply the missing ones
		sql, err := ioutil.ReadFile(file)
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

//	https://github.com/jackc/pgx/commit/aa8604b5c22989167e7158ecb1f6e7b8ddfebf04
//	if err = db.Ping(ctx); err != nil {
//		log.Fatal(err)
//	}

	// apply database migrations
	if err = migrate(ctx, db, "migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	// use postgres for session data
	session := scs.New()
	session.Store = pgxstore.New(db)

	store := NewPostgresStorage(ctx, db)
	app := &Application{
		session: session,
		store:   store,
	}

	if *addblog {
		feedURL := os.Args[2]
		siteURL := os.Args[3]
		fmt.Printf("adding blog: %s\n", feedURL)

		fp := gofeed.NewParser()
		blog, err := fp.ParseURL(feedURL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  found: %s\n", blog.Title)

		_, err = app.store.CreateBlog(feedURL, siteURL, blog.Title)
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

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(app.HandleIndex))
	mux.Handle("/about", http.HandlerFunc(app.HandleAbout))
	mux.Handle("/login", http.HandlerFunc(app.HandleLogin))
	mux.Handle("/register", http.HandlerFunc(app.HandleRegister))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

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

		// kick off hourly sync task
		go app.HourlySync()

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
		// local development setup
		fmt.Printf("Listening on %s\n", *addr)
		log.Fatal(http.ListenAndServe(*addr, mux))
	}
}
