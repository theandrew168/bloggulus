package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/*

Possible Errors (handle-able errors):

create (Exec) - constraint violation (conflict w/ cols)
	UUIDs will always be unique, but other fields could cause issues.
	For example, if adding a user with an email (NOT NULL UNIQUE) that
	already exists, this would cause a constraint violation.

list (CollectRows) - none
	Read only. No potential errors here: either returns some rows or none.

read (CollectOneRow) - does not exist
	Read only. Only potential error is not finding a row with the specified ID.

update (CollectOneRow) - does not exist, constraint violation (conflict w/ cols)
	CollectOneRow is used here to read the deleted record's updated_at and check for ErrNoRows.
	This is the most complex operation: multiple things could go wrong.
	  1. The record being updated doesn't exist (return ErrNotFound)
	  	 This could be a programming error on the caller's side: updating a record that doesn't exist yet.
	  2. The record being updated causes a constraint violation (dupe values in a UNIQUE column)
	  	 Probably need to communicate this back to the user in one way or another (return ErrConflict).
	  3. The record being updated was changed between fetch and update (TOCTOU race condition)
	     Based on Alex Edwards' approach to optimistic concurrency control in Let's Go Further.
		 The record exists, but was updated by someone (or something) else before the current
		 request completed. Probably need to tell the user to try again (return ErrNotFound).

delete (CollectOneRow) - does not exist
	CollectOneRow is used here to read the deleted record's ID and check for ErrNoRows.
	Only potential error is not finding a row with the specified ID. This could just
	ignore cases where the ID doesn't exist (and nothing gets deleted) but I think it is
	better UX / DX to _know_ if the delete was successful (204) vs no record was deleted (404).
	Could this also cause a constraint violation? For violating an FK or something?

*/

// TODO: add metadata to errors to make em more useful:
//   - what already exists
//   - what was missing
//   - what column(s) caused the conflict
var (
	ErrNotFound = errors.New("postgres: not found")
	ErrConflict = errors.New("postgres: conflict")

	// sentinel error used to rollback transactions
	ErrRollback = errors.New("postgres: rollback")
)

func CheckCreateError(err error) error {
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

func CheckListError(err error) error {
	return err
}

func CheckReadError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotFound
	default:
		return err
	}
}

func CheckUpdateError(err error) error {
	var pgErr *pgconn.PgError

	switch {
	// ErrNoRows in an update indicates a TOCTOU race condition (conflict)
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotFound
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

func CheckDeleteError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotFound
	default:
		return err
	}
}
