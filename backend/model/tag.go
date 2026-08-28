package model

import (
	"uuid"

	"github.com/theandrew168/bloggulus/backend/value"
)

type Tag struct {
	id   uuid.UUID
	name value.Name
	meta *Meta
}

func NewTag(name value.Name) (*Tag, error) {
	tag := Tag{
		id:   uuid.New(),
		name: name,
		meta: NewMeta(),
	}
	return &tag, nil
}

func LoadTag(id uuid.UUID, name value.Name, meta *Meta) *Tag {
	tag := Tag{
		id:   id,
		name: name,
		meta: meta,
	}
	return &tag
}

func (t *Tag) ID() uuid.UUID {
	return t.id
}

func (t *Tag) Name() value.Name {
	return t.name
}

func (t *Tag) Meta() *Meta {
	return t.meta
}
