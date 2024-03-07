package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/database"
	"github.com/theandrew168/bloggulus/backend/feed"
	"github.com/theandrew168/bloggulus/backend/migrate"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/task"
	"github.com/theandrew168/bloggulus/backend/web"
	"github.com/theandrew168/bloggulus/frontend"
)

//go:embed migrations
var migrationsFS embed.FS

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
	applied, err := migrate.Migrate(pool, migrationsFS)
	if err != nil {
		logger.Println(err.Error())
		return 1
	}

	for _, migration := range applied {
		logger.Printf("applied migration: %s\n", migration)
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
			if err == storage.ErrExist {
				logger.Println("  already exists")
			} else {
				logger.Println(err)
				return 1
			}
		}

		return 0
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// create a context that cancels upon receiving an interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	app := web.NewApplication(logger, store, frontend.Frontend)

	// let port be overridden by an env var
	port := cfg.Port
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	addr := fmt.Sprintf("127.0.0.1:%s", port)

	// start the web server in the background
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := app.Run(ctx, addr)
		if err != nil {
			logger.Println(err.Error())
		}
	}()

	// init task worker
	worker := task.NewWorker(logger)
	syncBlogs := worker.SyncBlogs(store, reader)

	// start worker in the background (standalone mode by default)
	wg.Add(1)
	go func() {
		defer wg.Done()

		// kick off blog sync task
		go syncBlogs.Run(1 * time.Hour)

		<-ctx.Done()

		logger.Println("stopping worker")
		worker.Wait()
		logger.Println("stopped worker")
	}()

	// wait for the web server and worker to stop
	wg.Wait()

	return 0
}
