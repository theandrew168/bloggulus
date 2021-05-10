package model

import (
	"context"
	"time"
)

type Session struct {
	SessionID string
	AccountID int
	Expiry    time.Time

	Account Account
}

type SessionStorage interface {
    Create(ctx context.Context, session *Session) (*Session, error)
    Read(ctx context.Context, sessionID string) (*Session, error)
    Delete(ctx context.Context, sessionID string) error
    DeleteExpired(ctx context.Context) error
}