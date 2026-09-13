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

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/config"
	webfeed "github.com/theandrew168/bloggulus/backend/feed/web"
	"github.com/theandrew168/bloggulus/backend/job"
	"github.com/theandrew168/bloggulus/backend/otel"
	"github.com/theandrew168/bloggulus/backend/postgres"
	webquery "github.com/theandrew168/bloggulus/backend/query/web"
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

	tracerCtx := context.Background()

	res := otel.NewResource()

	loggerProvider, err := otel.NewLogger(tracerCtx, res)
	if err != nil {
		return err
	}
	defer loggerProvider.Shutdown(tracerCtx)

	tracerProvider, err := otel.NewTracer(tracerCtx, res)
	if err != nil {
		return err
	}
	defer tracerProvider.Shutdown(tracerCtx)

	// Load the application's config file.
	conf, err := config.ReadFile(*configFilePath)
	if err != nil {
		return err
	}

	// Configure the database connection pool.
	poolConfig, err := postgres.PoolConfig(conf.DatabaseURI)
	if err != nil {
		return err
	}

	// Enable query-level tracing.
	poolConfig.ConnConfig.Tracer = otel.NewQueryTracer(tracerProvider)

	// Open a database connection pool.
	pool, err := postgres.ConnectPool(poolConfig)
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

	feedFetcher := webfeed.NewFeedFetcher()

	// Init the database storage interfaces.
	repo := repository.New(pool)
	cmd := command.New(repo, feedFetcher)
	qry := webquery.New(pool)

	// Init the sync service and do an initial sync.
	syncService := job.NewSyncService(cmd)

	// Init the session service and clear any expired session tokens.
	sessionService := job.NewSessionService(cmd)

	// Let systemd know that we are good to go (no-op if not using systemd).
	daemon.SdNotify(false, daemon.SdNotifyReady)

	handler := web.Handler(publicFS, conf, cmd, qry, syncService)

	// Let the web server port be overridden by an env var.
	port := "5000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	addr := fmt.Sprintf("127.0.0.1:%s", port)

	// Create a context that cancels upon receiving an interrupt signal.
	cancelCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	// Start the web server in the background.
	wg.Go(func() {
		err := web.Run(cancelCtx, handler, addr)
		if err != nil {
			slog.Error("error running web server",
				"error", err.Error(),
			)
		}
	})

	// Start the sync service in the background.
	wg.Go(func() {
		err := syncService.Run(cancelCtx)
		if err != nil {
			slog.Error("error running sync service",
				"error", err.Error(),
			)
		}
	})

	// Start the session cleanup service in the background.
	wg.Go(func() {
		err := sessionService.Run(cancelCtx)
		if err != nil {
			slog.Error("error running session service",
				"error", err.Error(),
			)
		}
	})

	// Wait for all services to stop.
	wg.Wait()

	return nil
}
