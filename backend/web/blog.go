package web

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleBlogRead(repo *repository.Repository) http.Handler {
	tmpl := page.NewBlog()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		blog, err := repo.Blog().Read(blogID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		posts, err := repo.Post().List(blog, 100, 0)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		data := page.BlogData{
			BaseData: util.TemplateBaseData(r, w),

			ID:       blog.ID(),
			Title:    blog.Title(),
			SiteURL:  blog.SiteURL(),
			FeedURL:  blog.FeedURL(),
			SyncedAt: blog.SyncedAt(),
		}
		for _, post := range posts {
			blogPost := page.PostData{
				ID:          post.ID(),
				BlogID:      post.BlogID(),
				Title:       post.Title(),
				PublishedAt: post.PublishedAt(),
			}
			data.Posts = append(data.Posts, blogPost)
		}

		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleBlogDeleteForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		blog, err := repo.Blog().Read(blogID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		err = repo.Blog().Delete(blog)
		if err != nil {
			util.DeleteErrorResponse(w, r, err)
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
