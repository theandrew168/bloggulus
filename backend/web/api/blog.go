package api

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

type jsonBlog struct {
	ID      uuid.UUID `json:"id"`
	FeedURL string    `json:"feedURL"`
	SiteURL string    `json:"siteURL"`
	Title   string    `json:"title"`
}

func marshalBlog(blog *admin.Blog) jsonBlog {
	b := jsonBlog{
		ID:      blog.ID(),
		FeedURL: blog.FeedURL(),
		SiteURL: blog.SiteURL(),
		Title:   blog.Title(),
	}
	return b
}

func (app *Application) handleBlogRead() http.HandlerFunc {
	type response struct {
		Blog jsonBlog `json:"blog"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()

		id, err := uuid.Parse(flow.Param(r.Context(), "id"))
		if err != nil {
			v.AddError("id", "must be a valid UUID")
			app.badRequestResponse(w, r, v.Errors)
			return
		}

		blog, err := app.storage.Blog().Read(id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		resp := response{
			Blog: marshalBlog(blog),
		}
		err = writeJSON(w, 200, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleBlogList() http.HandlerFunc {
	type response struct {
		Blogs []jsonBlog `json:"blogs"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		qs := r.URL.Query()

		limit := readInt(qs, "limit", 20, v)
		v.Check(limit >= 0, "limit", "must be positive")
		v.Check(limit <= 50, "limit", "must be less than or equal to 50")

		offset := readInt(qs, "offset", 0, v)
		v.Check(offset >= 0, "offset", "must be positive")

		if !v.Valid() {
			app.badRequestResponse(w, r, v.Errors)
			return
		}

		blogs, err := app.storage.Blog().List(limit, offset)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		resp := response{
			// use make here to encode JSON as "[]" instead of "null" if empty
			Blogs: make([]jsonBlog, 0),
		}

		for _, blog := range blogs {
			resp.Blogs = append(resp.Blogs, marshalBlog(blog))
		}

		err = writeJSON(w, 200, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}
