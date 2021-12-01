package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

func (app *Application) HandleReadBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	blog, err := app.storage.ReadBlog(context.Background(), id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, 200, envelope{"blog": blog}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleReadBlogs(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	blogs, err := app.storage.ReadBlogs(context.Background(), limit, offset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, 200, envelope{"blogs": blogs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
