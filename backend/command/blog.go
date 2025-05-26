package command

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
)

var ErrBlogNotFound = errors.New("blog not found")

func (cmd *Command) DeleteBlog(blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		blog, err := tx.Blog().Read(blogID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrBlogNotFound
			}

			return err
		}

		err = tx.Blog().Delete(blog)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrBlogNotFound
			}

			return err
		}

		slog.Info("blog deleted",
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		return nil
	})
}
