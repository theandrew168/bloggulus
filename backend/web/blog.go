package web

import (
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/command"
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

		posts, err := repo.Post().ListByBlog(blog)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		data := page.BlogData{
			BaseData: util.GetTemplateBaseData(r, w),

			Blog:  blog,
			Posts: posts,
		}
		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleBlogDeleteForm(cmd *command.Command) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		err = cmd.DeleteBlog(blogID)
		if err != nil {
			if errors.Is(err, command.ErrBlogNotFound) {
				util.NotFoundResponse(w, r)
				return
			}

			util.InternalServerErrorResponse(w, r, err)
			return
		}

		// Redirect back to the blogs page.
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}
