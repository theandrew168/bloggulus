package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
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

func HandleTagCreate(store *storage.Storage) http.Handler {
	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		Tag jsonTag `json:"tag"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := util.NewErrors()
		body := util.ReadBody(w, r)

		var req request
		err := util.ReadJSON(body, &req, true)
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		e.CheckField(req.Name != "", "must be provided", "name")

		if !e.Valid() {
			util.FailedValidationResponse(w, r, e)
			return
		}

		tag, err := model.NewTag(req.Name)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		err = store.Tag().Create(tag)
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
	})
}

func HandleTagList(store *storage.Storage) http.Handler {
	type response struct {
		Count int       `json:"count"`
		Tags  []jsonTag `json:"tags"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := util.NewErrors()
		qs := r.URL.Query()

		// check pagination params
		page := util.ReadInt(qs, "page", 1, e)
		e.CheckField(page >= 1, "Page must be greater than or equal to 1", "page")

		size := util.ReadInt(qs, "size", 20, e)
		e.CheckField(size >= 1, "Size must be greater than or equal to 1", "size")
		e.CheckField(size <= 50, "Size must be less than or equal to 50", "size")

		if !e.Valid() {
			util.FailedValidationResponse(w, r, e)
			return
		}

		limit, offset := util.PageSizeToLimitOffset(page, size)

		var count int
		var tags []*model.Tag

		var g errgroup.Group
		g.Go(func() error {
			var err error
			count, err = store.Tag().Count()
			return err
		})
		g.Go(func() error {
			var err error
			tags, err = store.Tag().List(limit, offset)
			return err
		})

		err := g.Wait()
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Count: count,
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
	})
}

func HandleTagDelete(store *storage.Storage) http.Handler {
	type response struct {
		Tag jsonTag `json:"tag"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tagID, err := uuid.Parse(r.PathValue("tagID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		tag, err := store.Tag().Read(tagID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		err = store.Tag().Delete(tag)
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
	})
}
