package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/timeutil"
)

const (
	// Clear out expired sessions every few minutes.
	ClearExpiredSessionsInterval = 5 * time.Minute
)

type SessionService struct {
	repo *repository.Repository
}

func NewSessionService(repo *repository.Repository) *SessionService {
	s := SessionService{
		repo: repo,
	}
	return &s
}

func (s *SessionService) Run(ctx context.Context) error {
	// Clear out any expired sessions at service startup.
	err := s.ClearExpiredSessions()
	if err != nil {
		slog.Error("error clearing expired sessions",
			"error", err.Error(),
		)
	}

	// Then run again every "internal" until stopped (by the context being canceled).
	ticker := time.NewTicker(SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping session service")
			slog.Info("stopped session service")
			return nil
		case <-ticker.C:
			err := s.ClearExpiredSessions()
			if err != nil {
				slog.Error("error clearing expired sessions",
					"error", err.Error(),
				)
			}
		}
	}
}

func (s *SessionService) ClearExpiredSessions() error {
	now := timeutil.Now()
	return s.repo.Session().DeleteExpired(now)
}
