package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

type jsonPost struct {
	ID          uuid.UUID `json:"id"`
	BlogID      uuid.UUID `json:"blogID"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	PublishedAt time.Time `json:"publishedAt"`
}

func marshalPost(post *model.Post) jsonPost {
	p := jsonPost{
		ID:          post.ID(),
		BlogID:      post.BlogID(),
		URL:         post.URL(),
		Title:       post.Title(),
		PublishedAt: post.PublishedAt(),
	}
	return p
}

func HandlePostRead(store *storage.Storage) http.Handler {
	type response struct {
		Post jsonPost `json:"post"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := uuid.Parse(r.PathValue("postID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		post, err := store.Post().Read(postID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

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
	})
}

func HandlePostList(store *storage.Storage) http.Handler {
	type response struct {
		Count int        `json:"count"`
		Posts []jsonPost `json:"posts"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		blog, err := store.Blog().Read(blogID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

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
		var posts []*model.Post

		var g errgroup.Group
		g.Go(func() error {
			var err error
			count, err = store.Post().Count(blog)
			return err
		})
		g.Go(func() error {
			var err error
			posts, err = store.Post().List(blog, limit, offset)
			return err
		})

		err = g.Wait()
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Count: count,
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
	})
}

func HandlePostDelete(store *storage.Storage) http.Handler {
	type response struct {
		Post jsonPost `json:"post"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := uuid.Parse(r.PathValue("postID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		post, err := store.Post().Read(postID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.NotFoundResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		err = store.Post().Delete(post)
		if err != nil {
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
	})
}
