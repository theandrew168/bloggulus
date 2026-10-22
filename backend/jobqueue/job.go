package jobqueue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"uuid"

	"github.com/riverqueue/river"
	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/otel"
	"github.com/theandrew168/bloggulus/backend/postgres"
	webquery "github.com/theandrew168/bloggulus/backend/query/web"
	"github.com/theandrew168/bloggulus/backend/timeutil"
	"github.com/theandrew168/bloggulus/backend/value"
)

type DeleteExpiredSessionsArgs struct{}

func (DeleteExpiredSessionsArgs) Kind() string {
	return "DeleteExpiredSessions"
}

type DeleteExpiredSessionsWorker struct {
	river.WorkerDefaults[DeleteExpiredSessionsArgs]

	cmd *command.Command
}

func (w *DeleteExpiredSessionsWorker) Work(ctx context.Context, job *river.Job[DeleteExpiredSessionsArgs]) error {
	tracerCtx, span := otel.GetTracer().Start(ctx, "Job_DeleteExpiredSessions")
	defer span.End()

	err := w.cmd.Auth().DeleteExpiredSessions(tracerCtx, timeutil.Now())
	if err != nil {
		slog.ErrorContext(tracerCtx, "error deleting expired sessions",
			"error", err.Error(),
		)

		return fmt.Errorf("error deleting expired sessions: %w", err)
	}

	return nil
}

type SyncBlogArgs struct {
	AccountID uuid.UUID
	FeedURL   string
}

func (SyncBlogArgs) Kind() string {
	return "SyncBlog"
}

type SyncBlogWorker struct {
	river.WorkerDefaults[SyncBlogArgs]

	cmd *command.Command
	qry *webquery.Query
}

func (w *SyncBlogWorker) Work(ctx context.Context, job *river.Job[SyncBlogArgs]) error {
	tracerCtx, span := otel.GetTracer().Start(ctx, "Job_SyncBlog")
	defer span.End()

	feedURL, err := value.NewURL(job.Args.FeedURL)
	if err != nil {
		slog.ErrorContext(tracerCtx, "error parsing feed URL",
			"error", err.Error(),
			"feedURL", job.Args.FeedURL,
		)

		return fmt.Errorf("error parsing feed URL: %w", err)
	}

	err = w.cmd.Sync().SyncBlog(tracerCtx, feedURL)
	if err != nil {
		slog.ErrorContext(tracerCtx, "error adding blog",
			"error", err.Error(),
			"feedURL", feedURL.Value(),
		)

		return fmt.Errorf("error adding blog: %w", err)
	}

	// The blog _should_ exist not if SyncBlog finished without errors.
	blog, err := w.qry.Blog().ReadDetailsByFeedURL(ctx, feedURL)
	if err != nil {
		slog.ErrorContext(tracerCtx, "error reading blog",
			"error", err.Error(),
			"feedURL", feedURL.Value(),
		)

		return fmt.Errorf("error reading blog: %w", err)
	}

	err = w.cmd.Account().FollowBlog(ctx, job.Args.AccountID, blog.ID)
	if err != nil {
		if !errors.Is(err, postgres.ErrConflict) {
			slog.ErrorContext(tracerCtx, "error following blog",
				"error", err.Error(),
				"account_id", job.Args.AccountID.String(),
				"blog_id", blog.ID.String(),
			)

			return fmt.Errorf("error following blog: %w", err)
		}
	}

	return nil
}

type SyncAllBlogsArgs struct{}

func (SyncAllBlogsArgs) Kind() string {
	return "SyncAllBlogs"
}

type SyncAllBlogsWorker struct {
	river.WorkerDefaults[SyncAllBlogsArgs]

	cmd *command.Command
}

func (w *SyncAllBlogsWorker) Work(ctx context.Context, job *river.Job[SyncAllBlogsArgs]) error {
	tracerCtx, span := otel.GetTracer().Start(ctx, "Job_SyncAllBlogs")
	defer span.End()

	err := w.cmd.Sync().SyncAllBlogs(tracerCtx)
	if err != nil {
		slog.ErrorContext(tracerCtx, "error syncing all blogs",
			"error", err.Error(),
		)

		return fmt.Errorf("error syncing all blogs: %w", err)
	}

	return nil
}
