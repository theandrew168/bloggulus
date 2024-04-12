package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

type jsonTag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func marshalTag(tag *admin.Tag) jsonTag {
	t := jsonTag{
		ID:   tag.ID(),
		Name: tag.Name(),
	}
	return t
}

func (app *Application) handleTagList() http.HandlerFunc {
	type response struct {
		Tags []jsonTag `json:"tags"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		qs := r.URL.Query()

		limit := readInt(qs, "limit", 20, v)
		v.Check(limit >= 0, "limit", "must be positive")
		v.Check(limit <= 50, "limit", "must be less than or equal to 50")

		offset := readInt(qs, "offset", 0, v)
		v.Check(offset >= 0, "offset", "must be positive")

		if !v.Valid() {
			util.BadRequestResponse(w, r, v.Errors)
			return
		}

		tags, err := app.storage.Tag().List(limit, offset)
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

		err = util.WriteJSON(w, 200, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}
