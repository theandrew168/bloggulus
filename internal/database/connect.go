package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
