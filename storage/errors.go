package storage

import (
	"errors"
)

var (
	ErrDuplicateModel = errors.New("storage: duplicate model")
	ErrNoModel        = errors.New("storage: no model")
)
