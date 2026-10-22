package jobqueue

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

func Migrate(pool *pgxpool.Pool) ([]string, error) {
	driver := riverpgxv5.New(pool)
	migrator, err := rivermigrate.New(driver, nil)
	if err != nil {
		return nil, err
	}

	res, err := migrator.Migrate(context.Background(), rivermigrate.DirectionUp, nil)
	if err != nil {
		return nil, err
	}

	var applied []string
	for _, version := range res.Versions {
		applied = append(applied, version.Name)
	}
	return applied, nil
}
