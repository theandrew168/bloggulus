package reader

import (
	"net/http"
	"time"

	"github.com/theandrew168/bloggulus/backend/domain/reader"
	"github.com/theandrew168/bloggulus/backend/web/api/util"
	"github.com/theandrew168/bloggulus/backend/web/api/validator"
)

type jsonPost struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	BlogTitle   string    `json:"blogTitle"`
	BlogURL     string    `json:"blogURL"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
}

func marshalPost(post *reader.Post) jsonPost {
	p := jsonPost{
		Title:       post.Title(),
		URL:         post.URL(),
		BlogTitle:   post.BlogTitle(),
		BlogURL:     post.BlogURL(),
		PublishedAt: post.PublishedAt(),
		Tags:        post.Tags(),
	}
	return p
}

func (app *Application) handlePostList() http.HandlerFunc {
	type response struct {
		Count int        `json:"count"`
		Posts []jsonPost `json:"posts"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		qs := r.URL.Query()

		// check search param
		q := qs.Get("q")

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

		var count int
		var posts []*reader.Post
		var err error

		if q != "" {
			count, err = app.store.Reader().Post().CountSearch(q)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}

			posts, err = app.store.Reader().Post().ListSearch(q, limit, offset)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}
		} else {
			count, err = app.store.Reader().Post().Count()
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}

			posts, err = app.store.Reader().Post().List(limit, offset)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}
		}

		resp := response{
			Count: count,
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
