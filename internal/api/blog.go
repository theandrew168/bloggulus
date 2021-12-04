package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/validator"
)

func (app *Application) HandleReadBlog(w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		v.AddError("id", "must be an integer")
		app.badRequestResponse(w, r, v.Errors)
		return
	}

	v.Check(id >= 0, "id", "must be positive")
	if !v.Valid() {
		app.badRequestResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	blog, err := app.storage.ReadBlog(ctx, id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
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

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	blogs, err := app.storage.ReadBlogs(ctx, limit, offset)
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
