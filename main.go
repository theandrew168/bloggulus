package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/domain/admin/service"
	"github.com/theandrew168/bloggulus/backend/postgres"
	adminStorage "github.com/theandrew168/bloggulus/backend/postgres/admin/storage"
	readerStorage "github.com/theandrew168/bloggulus/backend/postgres/reader/storage"
	"github.com/theandrew168/bloggulus/backend/web"
	"github.com/theandrew168/bloggulus/backend/web/fetch"
	"github.com/theandrew168/bloggulus/frontend"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	code := 0

	err := run()
	if err != nil {
		slog.Error(err.Error())
		code = 1
	}

	os.Exit(code)
}

func run() error {
	// set the program's local timezone to UTC
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}
	time.Local = utc

	// check for config file flag
	conf := flag.String("conf", "bloggulus.conf", "app config file")

	// check for action flags
	migrate := flag.Bool("migrate", false, "apply migrations and exit")
	addblog := flag.String("addblog", "", "rss / atom feed to add")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		return err
	}

	// open a database connection pool
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		return err
	}
	defer pool.Close()

	// apply database migrations
	applied, err := postgres.Migrate(pool, migrationsFS)
	if err != nil {
		return err
	}

	for _, migration := range applied {
		slog.Info("applied migration", "name", migration)
	}

	// exit now if just applying migrations
	if *migrate {
		return nil
	}

	// init database storage
	adminStore := adminStorage.New(pool)
	readerStore := readerStorage.New(pool)

	syncService := service.NewSyncService(adminStore, fetch.NewFeedFetcher(), fetch.NewPageFetcher())

	// add a blog and exit now if requested
	if *addblog != "" {
		feedURL := *addblog

		slog.Info("adding blog", "feedURL", feedURL)
		err = syncService.SyncBlog(feedURL)
		if err != nil {
			return err
		}

		return nil
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// create a context that cancels upon receiving an interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	app := web.NewApplication(frontend.Frontend, adminStore, readerStore)

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
			slog.Error(err.Error())
		}
	}()

	// wait for the web server and worker to stop
	wg.Wait()

	return nil
}
