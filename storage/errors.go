package storage

import (
	"errors"
)

var ErrDuplicateModel = errors.New("storage operation violates a unique constraint")
