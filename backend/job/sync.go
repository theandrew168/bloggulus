package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/theandrew168/bloggulus/backend/command"
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

func (s *SyncService) Run(ctx context.Context) error {
	// perform an initial sync at service startup
	err := s.cmd.Sync().SyncAllBlogs()
	if err != nil {
		slog.Error("error syncing blogs",
			"error", err.Error(),
		)
	}

	// then again every "internal" until stopped
	ticker := time.NewTicker(SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping sync service")
			slog.Info("stopped sync service")
			return nil
		case <-ticker.C:
			err := s.cmd.Sync().SyncAllBlogs()
			if err != nil {
				slog.Error("error syncing blogs",
					"error", err.Error(),
				)
			}
		}
	}
}
