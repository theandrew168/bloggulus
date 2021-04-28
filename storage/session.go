package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/model"
)

type Session interface {
	Create(ctx context.Context, session *model.Session) (*model.Session, error)
	Read(ctx context.Context, sessionID string) (*model.Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteExpired(ctx context.Context) error
}
