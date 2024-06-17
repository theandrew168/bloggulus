package admin

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	id   uuid.UUID
	name string

	createdAt time.Time
	updatedAt time.Time
}

func NewTag(name string) (*Tag, error) {
	now := time.Now().UTC()
	tag := Tag{
		id:   uuid.New(),
		name: name,

		createdAt: now,
		updatedAt: now,
	}
	return &tag, nil
}

func LoadTag(id uuid.UUID, name string, createdAt, updatedAt time.Time) *Tag {
	tag := Tag{
		id:   id,
		name: name,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &tag
}

func (t *Tag) ID() uuid.UUID {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Tag) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Tag) SetUpdatedAt(updatedAt time.Time) error {
	t.updatedAt = updatedAt
	return nil
}

func (t *Tag) CheckDelete() error {
	return nil
}
