package api

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

// TODO: Mat Ryer's req / resp closure pattern

func (app *Application) HandleReadBlog(w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	id, err := uuid.Parse(flow.Param(r.Context(), "id"))
	if err != nil {
		v.AddError("id", "must be a valid UUID")
		app.badRequestResponse(w, r, v.Errors)
		return
	}

	blog, err := app.storage.Blog.Read(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, 200, envelope{"blog": blog})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleReadBlogs(w http.ResponseWriter, r *http.Request) {
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

	blogs, err := app.storage.Blog.List(limit, offset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, 200, envelope{"blogs": blogs})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
