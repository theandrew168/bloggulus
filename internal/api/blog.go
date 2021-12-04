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
