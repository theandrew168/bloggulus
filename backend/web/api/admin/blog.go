package admin

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
	"github.com/theandrew168/bloggulus/backend/web/api/validator"
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

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			v.AddError("id", "must be a valid UUID")
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		blog, err := app.store.Admin().Blog().Read(id)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				util.NotFoundResponse(w, r)
				return
			}
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Blog: marshalBlog(blog),
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
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

		// check pagination params
		page := util.ReadInt(qs, "page", 1, v)
		v.Check(page >= 1, "page", "must be greater than or equal to 1")

		size := util.ReadInt(qs, "size", 20, v)
		v.Check(size >= 1, "size", "must be greater than or equal to 1")
		v.Check(size <= 50, "size", "must be less than or equal to 50")

		if !v.Valid() {
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		limit, offset := util.PageSizeToLimitOffset(page, size)

		blogs, err := app.store.Admin().Blog().List(limit, offset)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			// use make here to encode JSON as "[]" instead of "null" if empty
			Blogs: make([]jsonBlog, 0),
		}

		for _, blog := range blogs {
			resp.Blogs = append(resp.Blogs, marshalBlog(blog))
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}
