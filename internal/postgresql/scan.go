package postgresql

import (
	"github.com/jackc/pgx/v4"
)

func scan(row pgx.Row, dest ...interface{}) error {
	// TODO: pgerrcode.UniqueViolation -> core.ErrExist
	// TODO: pgx.ErrNoRows -> core.ErrNotExist
	// TODO: pgerrcode.AdminShutdown -> retry
	//	but how? need to call Query[Row] again?
	return nil
}

// TODO: instead write wrappers for?
// 	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
// 	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
