package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/otel"
)

const (
	// Check for new posts every SyncInterval.
	SyncInterval = 30 * time.Minute
)

type SyncService struct {
	cmd *command.Command
}

func NewSyncService(cmd *command.Command) *SyncService {
	s := SyncService{
		cmd: cmd,
	}
	return &s
}

func (s *SyncService) Run(cancelCtx context.Context) error {
	tracerCtx, span := otel.GetTracer().Start(context.Background(), "SyncService")

	// perform an initial sync at service startup
	err := s.cmd.Sync().SyncAllBlogs(tracerCtx)
	if err != nil {
		slog.Error("error syncing blogs",
			"error", err.Error(),
		)
	}

	span.End()

	// then again every "interval" until stopped
	ticker := time.NewTicker(SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cancelCtx.Done():
			slog.Info("stopping sync service")
			slog.Info("stopped sync service")
			return nil
		case <-ticker.C:
			tracerCtx, span := otel.GetTracer().Start(context.Background(), "SyncService")

			err := s.cmd.Sync().SyncAllBlogs(tracerCtx)
			if err != nil {
				slog.Error("error syncing blogs",
					"error", err.Error(),
				)
			}

			span.End()
		}
	}
}
