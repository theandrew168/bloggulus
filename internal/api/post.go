package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

func (app *Application) HandleReadPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}

	post, err := app.storage.ReadPost(context.Background(), id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			http.Error(w, "Not found", 404)
			return
		}
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	err = writeJSON(w, 200, envelope{"post": post}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
}

func (app *Application) HandleReadPosts(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	q := r.URL.Query().Get("q")

	var posts []core.Post
	if q != "" {
		// search if requested
		var err error
		posts, err = app.storage.SearchPosts(context.Background(), q, limit, offset)
		if err != nil {
			app.logger.Println(err)
			http.Error(w, "Internal server error", 500)
			return
		}
	} else {
		// else just read recent
		var err error
		posts, err = app.storage.ReadPosts(context.Background(), limit, offset)
		if err != nil {
			app.logger.Println(err)
			http.Error(w, "Internal server error", 500)
			return
		}
	}

	err = writeJSON(w, 200, envelope{"posts": posts}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
}
