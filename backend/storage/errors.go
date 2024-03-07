package storage

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/*

Possible Errors (handle-able errors):

create (exec) - constraint violation (conflict w/ cols)
	UUIDs will always be unique, but other fields could cause issues.
	For example, if adding a user with an email (NOT NULL UNIQUE) that
	already exists, this would cause a constraint violation.

list (collectRows) - none
	Read only. No potential errors here: either returns some rows or none.

read (collectOneRow) - does not exist
	Read only. Only potential error is not finding a row with the specified ID.

update (scan) - does not exist (coalesce to conflict), constraint violation (conflict w/ cols)
	Scan is used here to read the deleted record's version and check for ErrNoRows.
	This is the most complex operation: multiple things could go wrong.
	  1. The record being updated doesn't exist (indistinguishable from TOCTOU check, will appear as conflict)
	  	 This is technically a programming error on the caller's side: updating a record
		 that doesn't exist yet.
	  2. The record being updated causes a constraint violation (dupe values in a UNIQUE column)
	  	 Probably need to communicate this back to the user in one way or another.
	  3. The record being updated was changed between fetch and update (TOCTOU race condition)
	     Based on Alex Edwards' approach to optimistic concurrency control in Let's Go Further.
		 The record exists, but was updated by someone (or something) else before the current
		 request completed. Probably need to tell the user to try again.

delete (scan) - does not exist
	Scan is used here to read the deleted record's ID and check for ErrNoRows.
	Only potential error is not finding a row with the specified ID. This could just
	ignore cases where the ID doesn't exist (and nothing gets deleted) but I think it is
	better UX / DX to _know_ if the delete was successful (204) vs no record was deleted (404).

*/

// TODO: add metadata to errors to make em more useful:
//   - what already exists
//   - what was missing
//   - what column(s) caused the conflict
var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("storage: already exists")
	ErrNotExist = errors.New("storage: does not exist")

	// storage errors
	ErrRetry    = errors.New("storage: retry storage operation")
	ErrConflict = errors.New("storage: conflict in storage operation")
)

func scan(row pgx.Row, dest ...interface{}) error {
	err := row.Scan(dest...)
	if err != nil {
		// check for empty result (from QueryRow)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExist
		}

		// check for more specific errors
		// https://github.com/jackc/pgx/wiki/Error-Handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// check for duplicate primary keys
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrExist
			}
		}

		// otherwise bubble the error as-is
		return err
	}

	return nil
}

func checkCreateError(err error) error {
	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr):
		switch {
		case pgerrcode.IsIntegrityConstraintViolation(pgErr.Code):
			return ErrConflict
		default:
			return err
		}
	default:
		return err
	}
}

func checkListError(err error) error {
	return err
}

func checkReadError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotExist
	default:
		return err
	}
}

func checkUpdateError(err error) error {
	var pgErr *pgconn.PgError

	switch {
	// ErrNoRows in an update indicates a TOCTOU race condition (conflict)
	case errors.Is(err, pgx.ErrNoRows):
		return ErrConflict
	case errors.As(err, &pgErr):
		switch {
		case pgerrcode.IsIntegrityConstraintViolation(pgErr.Code):
			return ErrConflict
		default:
			return err
		}
	default:
		return err
	}
}

func checkDeleteError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotExist
	default:
		return err
	}
}
