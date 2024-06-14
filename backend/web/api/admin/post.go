package admin

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
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

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			v.AddError("id", "must be a valid UUID")
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		post, err := app.store.Admin().Post().Read(id)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				util.NotFoundResponse(w, r)
				return
			}
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Post: marshalPost(post),
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
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

		posts, err := app.store.Admin().Post().List(limit, offset)
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

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}
