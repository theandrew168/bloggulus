package model

import (
	"errors"
)

var (
	ErrExist    = errors.New("model: already exists")
	ErrNotExist = errors.New("model: does not exist")
	ErrInternal = errors.New("model: storage error")
)
