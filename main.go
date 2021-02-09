package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mmcdole/gofeed"
	"golang.org/x/crypto/acme/autocert"
//	"golang.org/x/crypto/bcrypt"
)

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	_ ,err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type Application struct {
	blog        *models.BlogStorage
	post        *models.PostStorage
	sourcedPost *models.SourcedPostStorage
}

type IndexData struct {
	Posts []*models.SourcedPost
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	posts, err := app.sourcedPost.ReadRecent(r.Context(), 20)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := &IndexData{
		Posts: posts,
	}

	err = ts.Execute(w, data)
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
		updated = time.Now().AddDate(0, -3, 0)
	}

	_, err := app.post.Create(context.Background(), blogID, post.Link, post.Title, updated)
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
	blogs, err := app.blog.ReadAll(context.Background())
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

func main() {
	sessionID, err := GenerateSessionID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sessionID)

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
	db, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//	https://github.com/jackc/pgx/commit/aa8604b5c22989167e7158ecb1f6e7b8ddfebf04
	//	if err = db.Ping(ctx); err != nil {
	//		log.Fatal(err)
	//	}

	// apply database migrations
	if err = models.Migrate(db, "migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	app := &Application{
		blog:        models.NewBlogStorage(db),
		post:        models.NewPostStorage(db),
		sourcedPost: models.NewSourcedPostStorage(db),
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

		_, err = app.blog.Create(context.Background(), feedURL, siteURL, blog.Title)
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
		go http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
