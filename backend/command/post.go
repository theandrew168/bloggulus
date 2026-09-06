package command

import (
	"errors"
	"log/slog"
	"uuid"

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

func (cmd *PostCommand) DeletePost(postID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		post, err := tx.Post().Read(postID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrPostNotFound
			}

			return err
		}

		err = tx.Post().Delete(post)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrPostNotFound
			}

			return err
		}

		slog.Info("post deleted",
			"post_id", post.ID(),
			"post_title", post.Title(),
		)

		return nil
	})
}
