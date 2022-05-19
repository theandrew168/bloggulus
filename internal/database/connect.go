package database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(uri string) (*pgx.Conn, error) {
	ctx := context.Background()

	// open a database connection pool
	conn, err := pgx.Connect(ctx, uri)
	if err != nil {
		return nil, err
	}

	// test connection to ensure all is well
	if err = conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, err
	}

	return conn, nil
}

func ConnectPool(uri string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// open a database connection pool
	pool, err := pgxpool.Connect(ctx, uri)
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
