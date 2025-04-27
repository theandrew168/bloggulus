package web

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/coreos/go-systemd/v22/activation"
)

func Run(ctx context.Context, handler http.Handler, addr string) error {
	srv := http.Server{
		Handler: handler,
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

	systemdSocketListeners, err := activation.Listeners()
	if err != nil {
		return err
	}

	var listener net.Listener
	if len(systemdSocketListeners) > 0 {
		slog.Info("using systemd socket activation")
		listener = systemdSocketListeners[0]
	} else {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}
	}

	slog.Info("starting web server", "addr", listener.Addr().String())

	// listen and serve forever
	// ignore http.ErrServerClosed (expected upon stop)
	err = srv.Serve(listener)
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
