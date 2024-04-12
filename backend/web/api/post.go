package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
	"github.com/theandrew168/bloggulus/backend/web/api/validator"
)

type jsonPost struct {
	ID          uuid.UUID `json:"id"`
	BlogID      uuid.UUID `json:"blogID"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	PublishedAt time.Time `json:"publishedAt"`
}

func marshalPost(post *admin.Post) jsonPost {
	p := jsonPost{
		ID:          post.ID(),
		BlogID:      post.BlogID(),
		URL:         post.URL(),
		Title:       post.Title(),
		PublishedAt: post.PublishedAt(),
	}
	return p
}

func (app *Application) handlePostRead() http.HandlerFunc {
	type response struct {
		Post jsonPost `json:"post"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()

		id, err := uuid.Parse(flow.Param(r.Context(), "id"))
		if err != nil {
			v.AddError("id", "must be a valid UUID")
			util.BadRequestResponse(w, r, v.Errors)
			return
		}

		post, err := app.storage.Post().Read(id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				util.NotFoundResponse(w, r)
				return
			}
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Post: marshalPost(post),
		}

		err = util.WriteJSON(w, 200, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handlePostList() http.HandlerFunc {
	type response struct {
		Posts []jsonPost `json:"posts"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		qs := r.URL.Query()

		limit := util.ReadInt(qs, "limit", 20, v)
		v.Check(limit >= 0, "limit", "must be positive")
		v.Check(limit <= 50, "limit", "must be less than or equal to 50")

		offset := util.ReadInt(qs, "offset", 0, v)
		v.Check(offset >= 0, "offset", "must be positive")

		if !v.Valid() {
			util.BadRequestResponse(w, r, v.Errors)
			return
		}

		posts, err := app.storage.Post().List(limit, offset)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			// use make here to encode JSON as "[]" instead of "null" if empty
			Posts: make([]jsonPost, 0),
		}

		for _, post := range posts {
			resp.Posts = append(resp.Posts, marshalPost(post))
		}

		err = util.WriteJSON(w, 200, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}
