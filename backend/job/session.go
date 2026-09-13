package job

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"

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

func (s *SessionService) Run(cancelCtx context.Context) error {
	tracerCtx, span := otel.Tracer("bloggulus").Start(context.Background(), "SessionService")

	// Clear out any expired sessions at service startup.
	err := s.cmd.Auth().DeleteExpiredSessions(tracerCtx, timeutil.Now())
	if err != nil {
		slog.Error("error clearing expired sessions",
			"error", err.Error(),
		)
	}

	span.End()

	// Then run again every "interval" until stopped (by the context being canceled).
	ticker := time.NewTicker(ClearExpiredSessionsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cancelCtx.Done():
			slog.Info("stopping session service")
			slog.Info("stopped session service")
			return nil
		case <-ticker.C:
			tracerCtx, span := otel.Tracer("bloggulus").Start(context.Background(), "SessionService")

			err := s.cmd.Auth().DeleteExpiredSessions(tracerCtx, timeutil.Now())
			if err != nil {
				slog.Error("error clearing expired sessions",
					"error", err.Error(),
				)
			}

			span.End()
		}
	}
}
