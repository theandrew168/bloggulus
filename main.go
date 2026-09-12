package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/coreos/go-systemd/v22/daemon"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/config"
	webfeed "github.com/theandrew168/bloggulus/backend/feed/web"
	"github.com/theandrew168/bloggulus/backend/job"
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

	// Create a context that cancels upon receiving an interrupt signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// TODO: What breaks if I remove this?
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	shutdownLogger, err := initOpenTelemetryLogger(ctx)
	if err != nil {
		return err
	}
	defer shutdownLogger()

	shutdownTracer, err := initOpenTelemetryTracer(ctx)
	if err != nil {
		return err
	}
	defer shutdownTracer()

	webHandler := web.Handler(publicFS, conf, cmd, qry, syncService)

	otelWebHandler := otelhttp.NewHandler(
		webHandler,
		"http-server-fallback", // The base "operation" name fallback
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			if r.Pattern != "" {
				// Returns exactly what matched, e.g., "GET /users/{id}"
				return r.Pattern
			}
			// Fallback for unmatched/404 routes
			return fmt.Sprintf("%s %s", r.Method, operation)
		}),
	)

	// Let the web server port be overridden by an env var.
	port := "5000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	addr := fmt.Sprintf("127.0.0.1:%s", port)

	var wg sync.WaitGroup

	// Start the web server in the background.
	wg.Go(func() {
		err := web.Run(ctx, otelWebHandler, addr)
		if err != nil {
			slog.Error("error running web server",
				"error", err.Error(),
			)
		}
	})

	// Start the sync service in the background.
	wg.Go(func() {
		err := syncService.Run(ctx)
		if err != nil {
			slog.Error("error running sync service",
				"error", err.Error(),
			)
		}
	})

	// Start the session cleanup service in the background.
	wg.Go(func() {
		err := sessionService.Run(ctx)
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

type ShutdownFunc func() error

func initOpenTelemetryLogger(ctx context.Context) (ShutdownFunc, error) {
	victoriaLogsExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpointURL("http://localhost:9428/insert/opentelemetry/v1/logs"),
	)
	if err != nil {
		return nil, err
	}

	// TODO: Use NewSimpleProcessor for dev
	// TODO: Use NewBatchProcessor for prod

	loggerProvider := log.NewLoggerProvider(log.WithProcessor(log.NewSimpleProcessor(victoriaLogsExporter)))
	global.SetLoggerProvider(loggerProvider)

	shutdown := func() error {
		return loggerProvider.Shutdown(ctx)
	}

	return shutdown, nil
}

func initOpenTelemetryTracer(ctx context.Context) (ShutdownFunc, error) {
	victoriaTracesExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL("http://localhost:10428/insert/opentelemetry/v1/traces"),
	)
	if err != nil {
		return nil, err
	}

	// TODO: Use WithSyncer for dev
	// TODO: Use WithBatcher for prod

	tracerProvider := trace.NewTracerProvider(trace.WithSyncer(victoriaTracesExporter))
	otel.SetTracerProvider(tracerProvider)

	shutdown := func() error {
		return tracerProvider.Shutdown(ctx)
	}

	return shutdown, nil
}
