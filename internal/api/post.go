package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/bloggulus"
	"github.com/theandrew168/bloggulus/internal/database"
	"github.com/theandrew168/bloggulus/internal/validator"
)

func (app *Application) HandleReadPost(w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
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

	post, err := app.storage.Post.Read(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
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

	var posts []bloggulus.Post
	if q != "" {
		// search if requested
		var err error
		posts, err = app.storage.Post.Search(q, limit, offset)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		// else just read recent
		var err error
		posts, err = app.storage.Post.ReadAll(limit, offset)
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
