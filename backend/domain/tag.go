package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTag(name string) Tag {
	now := time.Now()
	tag := Tag{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return tag
}
