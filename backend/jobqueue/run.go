package jobqueue

import (
	"context"
	"errors"
	"log/slog"
)

func Run(cancelCtx context.Context, queue *Queue) error {
	slog.Info("starting job queue")

	err := queue.Start(cancelCtx)
	if err != nil {
		return err
	}

	<-cancelCtx.Done()

	slog.Info("stopping job queue")
	err = queue.Stop(cancelCtx)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	slog.Info("stopped job queue")
	return nil
}
