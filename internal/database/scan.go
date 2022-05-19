package database

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

func Scan(row pgx.Row, dest ...interface{}) error {
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
			// check for stale connections (database restarted)
			if pgErr.Code == pgerrcode.AdminShutdown {
				return ErrRetry
			}
		}

		// otherwise bubble the error as-is
		return err
	}

	return nil
}
