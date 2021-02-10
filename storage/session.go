package storage

import (
	"context"

	"github.com/theandrew168/bloggulus/models"
)

type Session interface {
	Create(ctx context.Context, session *models.Session) (*models.Session, error)
	Read(ctx context.Context, sessionID string) (*models.Session, error)
	Delete(ctx context.Context, sessionID string) error
}
