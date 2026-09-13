package command

import (
	"context"
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

func (cmd *BlogCommand) DeleteBlog(ctx context.Context, blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		blog, err := tx.Blog().Read(ctx, blogID)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrBlogNotFound
			}

			return err
		}

		err = tx.Blog().Delete(ctx, blog)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return ErrBlogNotFound
			}

			return err
		}

		slog.Info("blog deleted",
			"blog_id", blog.ID().String(),
			"blog_title", blog.Title().Value(),
		)

		return nil
	})
}

func (cmd *BlogCommand) HideBlog(ctx context.Context, blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		blog, err := tx.Blog().Read(ctx, blogID)
		if err != nil {
			return err
		}

		err = blog.SetIsPublic(false)
		if err != nil {
			return err
		}

		err = tx.Blog().Update(ctx, blog)
		if err != nil {
			return err
		}

		slog.Info("blog hidden",
			"blog_id", blog.ID().String(),
			"blog_title", blog.Title().Value(),
		)

		return nil
	})
}

func (cmd *BlogCommand) ShowBlog(ctx context.Context, blogID uuid.UUID) error {
	return cmd.repo.WithTransaction(func(tx *repository.Repository) error {
		blog, err := tx.Blog().Read(ctx, blogID)
		if err != nil {
			return err
		}

		err = blog.SetIsPublic(true)
		if err != nil {
			return err
		}

		err = tx.Blog().Update(ctx, blog)
		if err != nil {
			return err
		}

		slog.Info("blog shown",
			"blog_id", blog.ID().String(),
			"blog_title", blog.Title().Value(),
		)

		return nil
	})
}
