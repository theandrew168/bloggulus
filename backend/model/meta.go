package model

import (
	"time"

	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/value"
)

type Meta struct {
	createdAt time.Time
	updatedAt time.Time
	version   value.Count
}

func NewMeta() *Meta {
	version, err := value.NewCount(1)
	if err != nil {
		panic(err)
	}

	m := Meta{
		createdAt: timeutil.Now(),
		updatedAt: timeutil.Now(),
		version:   version,
	}
	return &m
}

func LoadMeta(createdAt, updatedAt time.Time, version value.Count) *Meta {
	return &Meta{
		createdAt: createdAt,
		updatedAt: updatedAt,
		version:   version,
	}
}

func (m *Meta) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Meta) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *Meta) Version() value.Count {
	return m.version
}

func (m *Meta) Update(updatedAt time.Time) {
	m.updatedAt = updatedAt
}

func (m *Meta) IncrementVersion(updatedAt time.Time) {
	m.updatedAt = updatedAt
	m.version = m.version.Increment()
}
