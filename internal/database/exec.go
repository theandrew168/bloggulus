package database

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func Exec(db Conn, ctx context.Context, stmt string, args ...interface{}) error {
	_, err := db.Exec(ctx, stmt, args...)
	if err != nil {
		// check for more specific errors
		// https://github.com/jackc/pgx/wiki/Error-Handling
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
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
