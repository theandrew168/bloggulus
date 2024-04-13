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
		Posts []jsonPost `json:"posts"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		qs := r.URL.Query()

		// check search param
		q := qs.Get("q")

		limit := util.ReadInt(qs, "limit", 20, v)
		v.Check(limit >= 0, "limit", "must be positive")
		v.Check(limit <= 50, "limit", "must be less than or equal to 50")

		offset := util.ReadInt(qs, "offset", 0, v)
		v.Check(offset >= 0, "offset", "must be positive")

		if !v.Valid() {
			util.BadRequestResponse(w, r, v.Errors)
			return
		}

		var posts []*reader.Post
		var err error

		if q != "" {
			posts, err = app.store.Reader().Post().Search(q, limit, offset)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}
		} else {
			posts, err = app.store.Reader().Post().List(limit, offset)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}
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
