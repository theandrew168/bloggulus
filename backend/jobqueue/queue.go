package jobqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/theandrew168/bloggulus/backend/command"
	webquery "github.com/theandrew168/bloggulus/backend/query/web"
)

const (
	DeleteExpiredSessionsInterval = 5 * time.Minute
	SyncAllBlogsInterval          = 30 * time.Minute
)

type Queue struct {
	client *river.Client[pgx.Tx]
}

func NewQueue(pool *pgxpool.Pool, cmd *command.Command, qry *webquery.Query) (*Queue, error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &DeleteExpiredSessionsWorker{cmd: cmd})
	river.AddWorker(workers, &SyncBlogWorker{cmd: cmd, qry: qry})
	river.AddWorker(workers, &SyncAllBlogsWorker{cmd: cmd})

	periodicJobs := []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(DeleteExpiredSessionsInterval),
			func() (river.JobArgs, *river.InsertOpts) {
				return DeleteExpiredSessionsArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(SyncAllBlogsInterval),
			func() (river.JobArgs, *river.InsertOpts) {
				return SyncAllBlogsArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	}

	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		MaxAttempts:  1,
		JobTimeout:   5 * time.Minute,
		PeriodicJobs: periodicJobs,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 16},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, err
	}

	q := Queue{
		client: riverClient,
	}
	return &q, nil
}

func (q *Queue) Start(ctx context.Context) error {
	return q.client.Start(ctx)
}

func (q *Queue) Stop(ctx context.Context) error {
	return q.client.Stop(ctx)
}

func (q *Queue) Insert(ctx context.Context, args river.JobArgs) error {
	res, err := q.client.Insert(ctx, args, nil)
	if err != nil {
		return err
	}

	if res.UniqueSkippedAsDuplicate {
		return fmt.Errorf("job skipped as duplicate")
	}

	return nil
}
