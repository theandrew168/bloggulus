package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type storage struct {
	conn *pgxpool.Pool
}

func NewStorage(conn *pgxpool.Pool) core.Storage {
	s := storage{
		conn: conn,
	}
	return &s
}
