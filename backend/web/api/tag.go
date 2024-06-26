package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/web/util"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

type jsonTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func marshalTag(tag *model.Tag) jsonTag {
	t := jsonTag{
		ID:   tag.ID(),
		Name: tag.Name(),
	}
	return t
}

func (app *Application) handleTagCreate() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		Tag jsonTag `json:"tag"`
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

		v.Check(req.Name != "", "name", "must be provided")

		if !v.Valid() {
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		tag, err := model.NewTag(req.Name)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		err = app.store.Tag().Create(tag)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Tag: marshalTag(tag),
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleTagList() http.HandlerFunc {
	type response struct {
		Tags []jsonTag `json:"tags"`
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

		tags, err := app.store.Tag().List(limit, offset)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			// use make here to encode JSON as "[]" instead of "null" if empty
			Tags: make([]jsonTag, 0),
		}

		for _, tag := range tags {
			resp.Tags = append(resp.Tags, marshalTag(tag))
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleTagDelete() http.HandlerFunc {
	type response struct {
		Tag jsonTag `json:"tag"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := uuid.Parse(r.PathValue("tagID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		tag, err := app.store.Tag().Read(tagID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		err = app.store.Tag().Delete(tag)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Tag: marshalTag(tag),
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}
