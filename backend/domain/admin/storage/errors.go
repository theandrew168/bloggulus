package storage

import "errors"

// TODO: add metadata to errors to make em more useful:
//   - what already exists
//   - what was missing
//   - what column(s) caused the conflict
var (
	ErrNotFound = errors.New("storage: not found")
	ErrConflict = errors.New("storage: conflict")
)
