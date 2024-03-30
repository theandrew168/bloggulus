package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/bloggulus/backend/config"
	"github.com/theandrew168/bloggulus/backend/domain/admin/service"
	"github.com/theandrew168/bloggulus/backend/postgres"
	adminStorage "github.com/theandrew168/bloggulus/backend/postgres/admin/storage"
	"github.com/theandrew168/bloggulus/backend/web"
	"github.com/theandrew168/bloggulus/frontend"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	// check for config file flag
	conf := flag.String("conf", "bloggulus.conf", "app config file")

	// check for action flags
	migrate := flag.Bool("migrate", false, "apply migrations and exit")
	addblog := flag.String("addblog", "", "rss / atom feed to add")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		// TODO: log this
		fmt.Println(err)
		return 1
	}

	// open a database connection pool
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		// TODO: log this
		fmt.Println(err)
		return 1
	}
	defer pool.Close()

	// apply database migrations
	applied, err := postgres.Migrate(pool, migrationsFS)
	if err != nil {
		// TODO: log this
		fmt.Println(err)
		return 1
	}

	for _, migration := range applied {
		// TODO: log this
		fmt.Printf("applied migration: %s\n", migration)
	}

	// exit now if just applying migrations
	if *migrate {
		return 0
	}

	// init database storage
	store := adminStorage.New(pool)

	syncService := service.NewSyncService(store, web.NewFeedFetcher(), web.NewPageFetcher())

	// add a blog and exit now if requested
	if *addblog != "" {
		feedURL := *addblog

		// TODO: log this
		fmt.Printf("adding blog: %s\n", feedURL)
		err = syncService.SyncBlog(feedURL)
		if err != nil {
			// TODO: log this
			fmt.Println(err)
			return 1
		}

		return 0
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// create a context that cancels upon receiving an interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	app := web.NewApplication(frontend.Frontend, store)

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
			// TODO: log this
			fmt.Println(err.Error())
		}
	}()

	// wait for the web server and worker to stop
	wg.Wait()

	return 0
}
