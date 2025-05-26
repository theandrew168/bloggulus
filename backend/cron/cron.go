package cron

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/theandrew168/bloggulus/backend/command"
)

type Cron struct {
	cmd *command.Command
}

type Job struct {
	Interval time.Duration
	RunNow   bool
	Run      func(cmd *command.Command) error
}

func NewCron(cmd *command.Command) *Cron {
	c := Cron{
		cmd: cmd,
	}
	return &c
}
func (c *Cron) Run(ctx context.Context, jobs []Job) error {
	var wg sync.WaitGroup

	for _, job := range jobs {
		// If the job is set to run immediately, run it now.
		if job.RunNow {
			err := job.Run(c.cmd)
			if err != nil {
				slog.Error("error running job",
					"error", err.Error(),
				)
			}
		}

		// Start the job in a goroutine that runs on the specified interval.
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()

			ticker := time.NewTicker(j.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					err := j.Run(c.cmd)
					if err != nil {
						slog.Error("error running job",
							"error", err.Error(),
						)
					}
				}
			}
		}(job)
	}

	// Block until all jobs are done.
	wg.Wait()
	return nil
}
