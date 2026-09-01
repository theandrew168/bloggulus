package model

import (
	"time"

	"github.com/theandrew168/bloggulus/backend/timeutil"
)

type Meta struct {
	createdAt time.Time
	updatedAt time.Time
}

func NewMeta() *Meta {
	m := Meta{
		createdAt: timeutil.Now(),
		updatedAt: timeutil.Now(),
	}
	return &m
}

type LoadMetaParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func LoadMeta(params LoadMetaParams) *Meta {
	return &Meta{
		createdAt: params.CreatedAt,
		updatedAt: params.UpdatedAt,
	}
}

func (m *Meta) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Meta) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *Meta) Update(now time.Time) {
	m.updatedAt = now
}
