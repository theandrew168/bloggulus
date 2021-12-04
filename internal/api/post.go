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

func (app *Application) HandleReadPost(w http.ResponseWriter, r *http.Request) {
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

	post, err := app.storage.ReadPost(ctx, id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, 200, envelope{"post": post})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleReadPosts(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	qs := r.URL.Query()

	limit := readInt(qs, "limit", 20, v)
	v.Check(limit >= 0, "limit", "must be positive")
	v.Check(limit <= 50, "limit", "must be less than or equal to 50")

	offset := readInt(qs, "offset", 0, v)
	v.Check(offset >= 0, "offset", "must be positive")

	q := qs.Get("q")

	if !v.Valid() {
		app.badRequestResponse(w, r, v.Errors)
		return
	}

	var posts []core.Post
	if q != "" {
		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		// search if requested
		var err error
		posts, err = app.storage.SearchPosts(ctx, q, limit, offset)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		// else just read recent
		var err error
		posts, err = app.storage.ReadPosts(ctx, limit, offset)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err := writeJSON(w, 200, envelope{"posts": posts})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
