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

type NewTagParams struct {
	Name value.Name
}

func NewTag(params NewTagParams) (*Tag, error) {
	tag := Tag{
		id:   uuid.New(),
		name: params.Name,
		meta: NewMeta(),
	}
	return &tag, nil
}

type LoadTagParams struct {
	ID   uuid.UUID
	Name value.Name
	Meta *Meta
}

func LoadTag(params LoadTagParams) *Tag {
	tag := Tag{
		id:   params.ID,
		name: params.Name,
		meta: params.Meta,
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
