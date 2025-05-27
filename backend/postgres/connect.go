package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Common interface for pgx.Conn, pgx.Pool, pgx.Tx, etc
// https://github.com/jackc/pgx/issues/644
type Conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Connect(uri string) (*pgx.Conn, error) {
	ctx := context.Background()

	config, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Ensure timestamps read from a timestamptz column retain their UTC location.
	// https://github.com/jackc/pgx/issues/1195#issuecomment-2002079265
	conn.TypeMap().RegisterType(&pgtype.Type{
		Name:  "timestamptz",
		OID:   pgtype.TimestamptzOID,
		Codec: &pgtype.TimestamptzCodec{ScanLocation: time.UTC},
	})

	// test connection to ensure all is well
	if err = conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, err
	}

	return conn, nil
}

func ConnectPool(uri string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	// Ensure timestamps read from a timestamptz column retain their UTC location.
	// https://github.com/jackc/pgx/issues/1195#issuecomment-2002079265
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "timestamptz",
			OID:   pgtype.TimestamptzOID,
			Codec: &pgtype.TimestamptzCodec{ScanLocation: time.UTC},
		})

		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// test connection to ensure all is well
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
