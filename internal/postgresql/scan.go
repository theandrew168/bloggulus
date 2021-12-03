package postgresql

import (
	"github.com/jackc/pgx/v4"
)

func scan(row pgx.Row, dest ...interface{}) error {
	// TODO: pgerrcode.UniqueViolation -> core.ErrExist
	// TODO: pgx.ErrNoRows -> core.ErrNotExist
	// TODO: pgerrcode.AdminShutdown -> retry
	//	but how? need to call Query[Row] again?
	//	just return ErrRetry or something like that?
	//	and make every caller go recursive if they see it?
	return nil
}
