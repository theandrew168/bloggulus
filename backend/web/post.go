package web

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandlePostRead(repo *repository.Repository) http.Handler {
	tmpl := page.NewPost()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := uuid.Parse(r.PathValue("postID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		post, err := repo.Post().Read(postID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		data := page.PostData{
			BaseData: util.GetTemplateBaseData(r, w),

			Post: post,
		}
		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandlePostDeleteForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		postID, err := uuid.Parse(r.PathValue("postID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		post, err := repo.Post().Read(postID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		err = repo.Post().Delete(post)
		if err != nil {
			util.DeleteErrorResponse(w, r, err)
			return
		}

		slog.Info("post deleted",
			"post_id", post.ID(),
			"post_title", post.Title(),
		)

		// Redirect back to the blog page for this post's blog.
		http.Redirect(w, r, fmt.Sprintf("/blogs/%s", blogID), http.StatusSeeOther)
	})
}
