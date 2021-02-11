package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/theandrew168/bloggulus/handlers"
	"github.com/theandrew168/bloggulus/models"
	"github.com/theandrew168/bloggulus/storage/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mmcdole/gofeed"
	"golang.org/x/crypto/acme/autocert"
	//	"golang.org/x/crypto/bcrypt"
)

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
	if err = postgres.Migrate(db, context.Background(), "migrations/*.sql"); err != nil {
		log.Fatal(err)
	}

	app := &handlers.Application{
		Account:     postgres.NewAccountStorage(db),
		Blog:        postgres.NewBlogStorage(db),
		Post:        postgres.NewPostStorage(db),
		Session:     postgres.NewSessionStorage(db),
		SourcedPost: postgres.NewSourcedPostStorage(db),
	}

	if *addblog {
		feedURL := os.Args[2]
		siteURL := os.Args[3]
		fmt.Printf("adding blog: %s\n", feedURL)

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  found: %s\n", feed.Title)

		blog := &models.Blog{
			FeedURL: feedURL,
			SiteURL: siteURL,
			Title:   feed.Title,
		}
		_, err = app.Blog.Create(context.Background(), blog)
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
	mux.Handle("/logout", http.HandlerFunc(app.HandleLogout))
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
//		go app.HourlySync()

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
