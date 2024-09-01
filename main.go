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

	"github.com/coreos/go-systemd/v22/daemon"

	"github.com/theandrew168/bloggulus/backend/config"
	fetch "github.com/theandrew168/bloggulus/backend/fetch/web"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/query"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web"
)

//go:embed public
var publicFS embed.FS

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
	// check for config file flag
	conf := flag.String("conf", "bloggulus.conf", "app config file")

	// check for action flags
	migrate := flag.Bool("migrate", false, "apply migrations and exit")
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
	store := storage.New(pool)
	query := query.New(pool)

	// init the sync service and do an initial sync
	syncService := service.NewSyncService(store, fetch.NewFeedFetcher(), fetch.NewPageFetcher())

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// create a context that cancels upon receiving an interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	webHandler := web.Handler(publicFS, store, query, syncService)

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

		err := web.Run(ctx, webHandler, addr)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	// start the sync service in the background
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := syncService.Run(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	// wait for the web server and sync service to stop
	wg.Wait()

	return nil
}
