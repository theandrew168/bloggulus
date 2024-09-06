package web

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/page"
)

func HandleBlogRead(repo *repository.Repository) http.Handler {
	tmpl := page.NewBlog()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			http.Error(w, "Not Found", 404)
			return
		}

		blog, err := repo.Blog().Read(blogID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Error(w, "Not Found", 404)
			default:
				http.Error(w, err.Error(), 500)
			}

			return
		}

		posts, err := repo.Post().List(blog, 20, 0)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := page.BlogData{
			ID:       blog.ID(),
			Title:    blog.Title(),
			SiteURL:  blog.SiteURL(),
			FeedURL:  blog.FeedURL(),
			SyncedAt: blog.SyncedAt(),
		}
		for _, post := range posts {
			blogPost := page.BlogDataPost{
				ID:          post.ID(),
				BlogID:      post.BlogID(),
				Title:       post.Title(),
				PublishedAt: post.PublishedAt(),
			}
			data.Posts = append(data.Posts, blogPost)
		}

		err = tmpl.Render(w, data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
}

func HandleBlogDeleteForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			http.Error(w, "Not Found", 404)
			return
		}

		blog, err := repo.Blog().Read(blogID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Error(w, "Not Found", 404)
			default:
				http.Error(w, err.Error(), 500)
			}

			return
		}

		err = repo.Blog().Delete(blog)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Error(w, "Not Found", 404)
			default:
				http.Error(w, err.Error(), 500)
			}

			return
		}

		slog.Info("blog deleted",
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		// Redirect back to the blogs page.
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}
