package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

func (app *Application) Run(ctx context.Context, addr string) error {
	srv := http.Server{
		Addr:    addr,
		Handler: app.Router(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// start a goro to watch for stop signal (context cancelled)
	stopError := make(chan error)
	go func() {
		<-ctx.Done()

		// give the web server 5 seconds to shutdown gracefully
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Info("stopping web server")

		// disable keepalives and shutdown gracefully
		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(timeout)
		if err != nil {
			stopError <- err
		}

		close(stopError)
	}()

	slog.Info("starting web server", "addr", srv.Addr)

	// listen and serve forever
	// ignore http.ErrServerClosed (expected upon stop)
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// check for errors that arose while stopping
	err = <-stopError
	if err != nil {
		return err
	}

	slog.Info("stopped web server")
	return nil
}
