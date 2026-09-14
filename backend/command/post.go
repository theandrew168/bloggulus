package command

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/otel"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
)

var ErrPostNotFound = errors.New("post: not found")

type PostCommand struct {
	repo *repository.Repository
}

func NewPost(repo *repository.Repository) *PostCommand {
	cmd := PostCommand{
		repo: repo,
	}
	return &cmd
}

func (cmd *PostCommand) DeletePost(ctx context.Context, postID uuid.UUID) error {
	ctx, span := otel.GetTracer().Start(ctx, "Command_Post_DeletePost")
	defer span.End()

	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		post, err := tx.Post().Read(ctx, postID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrPostNotFound
			}

			return err
		}

		err = tx.Post().Delete(ctx, post)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrPostNotFound
			}

			return err
		}

		slog.InfoContext(ctx, "post deleted",
			"post_id", post.ID().String(),
			"post_title", post.Title().Value(),
		)

		return nil
	})
}
