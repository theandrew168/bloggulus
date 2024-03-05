package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/bloggulus/backend/app"
	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/migrate"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/task"
)

//go:embed all:build
var buildFS embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	// log everything to stdout, include file name and line number
	logger := log.New(os.Stdout, "", 0)

	// check for config file flag
	conf := flag.String("conf", "bloggulus.conf", "app config file")

	// check for action flags
	migrateOnly := flag.Bool("migrate", false, "apply migrations and exit")
	addblog := flag.String("addblog", "", "rss / atom feed to add")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Println(err)
		return 1
	}

	// open a database connection pool
	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.Println(err)
		return 1
	}
	defer pool.Close()

	// apply database migrations
	if err = migrate.Migrate(pool, logger); err != nil {
		logger.Println(err)
		return 1
	}

	// exit now if just applying migrations
	if *migrateOnly {
		return 0
	}

	// init database storage
	store := storage.New(pool)

	// init default feed reader
	reader := feed.NewReader(logger)

	// add a blog and exit now if requested
	if *addblog != "" {
		feedURL := *addblog
		logger.Printf("adding blog: %s\n", feedURL)

		blog, err := reader.ReadBlog(feedURL)
		if err != nil {
			logger.Println(err)
			return 1
		}
		logger.Printf("  found: %s\n", blog.Title)

		err = store.Blog.Create(&blog)
		if err != nil {
			if err == database.ErrExist {
				logger.Println("  already exists")
			} else {
				logger.Println(err)
				return 1
			}
		}

		return 0
	}

	// init task worker
	worker := task.NewWorker(logger)

	// kick off blog sync task
	syncBlogs := worker.SyncBlogs(store, reader)
	go syncBlogs.Run(1 * time.Hour)

	// init main web handler
	addr := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	handler := app.New(logger, store, buildFS)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalln(err)
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	logger.Printf("started server on %s\n", addr)

	// kick off a goroutine to listen for SIGINT and SIGTERM
	shutdownError := make(chan error)
	go func() {
		// idle until a signal is caught
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// wait for background tasks to finish (no timeout here)
		logger.Println("stopping worker")
		worker.Wait()
		logger.Println("stopped worker")

		// give the web server 5 seconds to shutdown gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown the web server and track any errors
		logger.Println("stopping server")
		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	// serve the app, check for ErrServerClosed (expected after shutdown)
	err = srv.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Println(err)
		return 1
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		logger.Println(err)
		return 1
	}

	logger.Println("stopped server")
	return 0
}
