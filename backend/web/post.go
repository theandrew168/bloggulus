package web

import (
	"fmt"
	"io"
	"net/http"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/command"
	webquery "github.com/theandrew168/bloggulus/backend/query/web"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandlePostRead(qry *webquery.Query) http.Handler {
	tmpl := page.NewPost()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := uuid.Parse(r.PathValue("postID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		post, err := qry.Post().ReadDetailsByID(postID)
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

func HandlePostDeleteForm(cmd *command.Command) http.Handler {
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

		err = cmd.Post().DeletePost(postID)
		if err != nil {
			util.DeleteErrorResponse(w, r, err)
			return
		}

		// Redirect back to the blog page for this post's blog.
		http.Redirect(w, r, fmt.Sprintf("/blogs/%s", blogID), http.StatusSeeOther)
	})
}
