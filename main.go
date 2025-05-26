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
	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/job"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
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
		slog.Error("error running application",
			"error", err.Error(),
		)
		code = 1
	}

	os.Exit(code)
}

func run() error {
	// Check for the config file path flag.
	configFilePath := flag.String("conf", "bloggulus.conf", "app config file")

	// Check for any specific action flags.
	migrate := flag.Bool("migrate", false, "apply migrations and exit")
	flag.Parse()

	// Load the application's config file.
	conf, err := config.ReadFile(*configFilePath)
	if err != nil {
		return err
	}

	// Open a database connection pool.
	pool, err := postgres.ConnectPool(conf.DatabaseURI)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Apply any pending database migrations.
	applied, err := postgres.Migrate(pool, migrationsFS)
	if err != nil {
		return err
	}

	for _, migration := range applied {
		slog.Info("applied migration", "name", migration)
	}

	// Exit now if just applying migrations.
	if *migrate {
		return nil
	}

	// Init the database storage interfaces.
	repo := repository.New(pool)
	find := finder.New(pool)

	// Init the sync service and do an initial sync.
	feedFetcher := fetch.NewFeedFetcher()
	syncService := job.NewSyncService(repo, feedFetcher)

	// Init the session service and clear any expired session tokens.
	sessionService := job.NewSessionService(repo)

	// Let systemd know that we are good to go (no-op if not using systemd).
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// Create a context that cancels upon receiving an interrupt signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	webHandler := web.Handler(publicFS, conf, repo, find, syncService)

	// Let the web server port be overridden by an env var.
	port := conf.Port
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	addr := fmt.Sprintf("127.0.0.1:%s", port)

	// Start the web server in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := web.Run(ctx, webHandler, addr)
		if err != nil {
			slog.Error("error running web server",
				"error", err.Error(),
			)
		}
	}()

	// Start the sync service in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := syncService.Run(ctx)
		if err != nil {
			slog.Error("error running sync service",
				"error", err.Error(),
			)
		}
	}()

	// Start the session cleanup service in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := sessionService.Run(ctx)
		if err != nil {
			slog.Error("error running session service",
				"error", err.Error(),
			)
		}
	}()

	// Wait for all services to stop.
	wg.Wait()

	return nil
}
