package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/web/util"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

type jsonBlog struct {
	ID      uuid.UUID `json:"id"`
	FeedURL string    `json:"feedURL"`
	SiteURL string    `json:"siteURL"`
	Title   string    `json:"title"`
}

func marshalBlog(blog *model.Blog) jsonBlog {
	b := jsonBlog{
		ID:      blog.ID(),
		FeedURL: blog.FeedURL(),
		SiteURL: blog.SiteURL(),
		Title:   blog.Title(),
	}
	return b
}

// TODO: This should be async if the blog is new. If it is,
// just run it in the background and keep track of which user
// submitted it (to link it once complete). Should I invest
// in a proper queue + worker system? Probably River?
func (app *Application) handleBlogCreate() http.HandlerFunc {
	type request struct {
		FeedURL string `json:"feedURL"`
	}
	type response struct {
		Blog jsonBlog `json:"blog"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		body := util.ReadBody(w, r)

		var req request
		err := util.ReadJSON(body, &req, true)
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		v.Check(req.FeedURL != "", "feedURL", "must be provided")

		if !v.Valid() {
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		// Check if the blog already exists. If it does, return its details.
		blog, err := app.store.Blog().ReadByFeedURL(req.FeedURL)
		if err == nil {
			resp := response{
				Blog: marshalBlog(blog),
			}

			code := http.StatusOK
			err = util.WriteJSON(w, code, resp, nil)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}

			return
		}

		// At this point, the only "expected" error is ErrNotFound.
		if !errors.Is(err, postgres.ErrNotFound) {
			util.ServerErrorResponse(w, r, err)
			return
		}

		// Use the SyncService to add the new blog.
		blog, err = app.syncService.SyncBlog(req.FeedURL)
		if err != nil {
			switch {
			case errors.Is(err, fetch.ErrUnreachableFeed):
				v.AddError("feedURL", "must link to a valid feed")
				util.FailedValidationResponse(w, r, v.Errors())
			default:
				util.ServerErrorResponse(w, r, err)
			}
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

		blog, err := app.store.Blog().Read(id)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

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

		blogs, err := app.store.Blog().List(limit, offset)
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

func (app *Application) handleBlogDelete() http.HandlerFunc {
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

		blog, err := app.store.Blog().Read(id)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		err = app.store.Blog().Delete(blog)
		if err != nil {
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
