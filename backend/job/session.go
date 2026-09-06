package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

const (
	// Clear out expired sessions every few minutes.
	ClearExpiredSessionsInterval = 5 * time.Minute
)

type SessionService struct {
	cmd *command.Command
}

func NewSessionService(cmd *command.Command) *SessionService {
	s := SessionService{
		cmd: cmd,
	}
	return &s
}

func (s *SessionService) Run(ctx context.Context) error {
	// Clear out any expired sessions at service startup.
	err := s.cmd.Auth().DeleteExpiredSessions(timeutil.Now())
	if err != nil {
		slog.Error("error clearing expired sessions",
			"error", err.Error(),
		)
	}

	// Then run again every "internal" until stopped (by the context being canceled).
	ticker := time.NewTicker(ClearExpiredSessionsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping session service")
			slog.Info("stopped session service")
			return nil
		case <-ticker.C:
			err := s.cmd.Auth().DeleteExpiredSessions(timeutil.Now())
			if err != nil {
				slog.Error("error clearing expired sessions",
					"error", err.Error(),
				)
			}
		}
	}
}
