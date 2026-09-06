package command

import (
	"errors"
	"log/slog"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
)

var ErrBlogNotFound = errors.New("blog: not found")

type BlogCommand struct {
	repo *repository.Repository
}

func NewBlog(repo *repository.Repository) *BlogCommand {
	cmd := BlogCommand{
		repo: repo,
	}
	return &cmd
}

func (cmd *BlogCommand) DeleteBlog(blogID uuid.UUID) error {
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

func (cmd *BlogCommand) HideBlog(blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		blog, err := tx.Blog().Read(blogID)
		if err != nil {
			return err
		}

		err = blog.SetIsPublic(false)
		if err != nil {
			return err
		}

		err = tx.Blog().Update(blog)
		if err != nil {
			return err
		}

		slog.Info("blog hidden",
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		return nil
	})
}

func (cmd *BlogCommand) ShowBlog(blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		blog, err := tx.Blog().Read(blogID)
		if err != nil {
			return err
		}

		err = blog.SetIsPublic(true)
		if err != nil {
			return err
		}

		err = tx.Blog().Update(blog)
		if err != nil {
			return err
		}

		slog.Info("blog shown",
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		return nil
	})
}
