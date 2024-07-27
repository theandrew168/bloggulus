package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

type jsonBlog struct {
	ID       uuid.UUID `json:"id"`
	FeedURL  string    `json:"feedURL"`
	SiteURL  string    `json:"siteURL"`
	Title    string    `json:"title"`
	SyncedAt time.Time `json:"syncedAt"`
}

func marshalBlog(blog *model.Blog) jsonBlog {
	b := jsonBlog{
		ID:       blog.ID(),
		FeedURL:  blog.FeedURL(),
		SiteURL:  blog.SiteURL(),
		Title:    blog.Title(),
		SyncedAt: blog.SyncedAt(),
	}
	return b
}

// TODO: This should be async if the blog is new. If it is,
// just run it in the background and keep track of which user
// submitted it (to link it once complete). Should I invest
// in a proper queue + worker system? Probably River?
func HandleBlogCreate(store *storage.Storage, syncService *service.SyncService) http.Handler {
	type request struct {
		FeedURL string `json:"feedURL"`
	}
	type response struct {
		Blog jsonBlog `json:"blog"`
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

		e.CheckField(req.FeedURL != "", "must be provided", "feedURL")

		if !e.Valid() {
			util.FailedValidationResponse(w, r, e)
			return
		}

		// Check if the blog already exists. If it does, return its details.
		blog, err := store.Blog().ReadByFeedURL(req.FeedURL)
		if err == nil {
			resp := response{
				Blog: marshalBlog(blog),
			}

			code := http.StatusOK
			err = util.WriteJSON(w, code, resp, nil)
			if err != nil {
				util.ServerErrorResponse(w, r, err)
				return
			}

			return
		}

		// At this point, the only "expected" error is ErrNotFound.
		if !errors.Is(err, postgres.ErrNotFound) {
			util.ServerErrorResponse(w, r, err)
			return
		}

		// Use the SyncService to add the new blog.
		blog, err = syncService.SyncBlog(req.FeedURL)
		if err != nil {
			switch {
			case errors.Is(err, fetch.ErrUnreachableFeed):
				e.AddField("must link to a valid feed", "feedURL")
				util.FailedValidationResponse(w, r, e)
			default:
				util.ServerErrorResponse(w, r, err)
			}
			return
		}

		resp := response{
			Blog: marshalBlog(blog),
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	})
}

func HandleBlogRead(store *storage.Storage) http.Handler {
	type response struct {
		Blog jsonBlog `json:"blog"`
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

		resp := response{
			Blog: marshalBlog(blog),
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	})
}

func HandleBlogList(store *storage.Storage) http.Handler {
	type response struct {
		Count int        `json:"count"`
		Blogs []jsonBlog `json:"blogs"`
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
		var blogs []*model.Blog

		var g errgroup.Group
		g.Go(func() error {
			var err error
			count, err = store.Blog().Count()
			return err
		})
		g.Go(func() error {
			var err error
			blogs, err = store.Blog().List(limit, offset)
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
			Blogs: make([]jsonBlog, 0),
		}

		for _, blog := range blogs {
			resp.Blogs = append(resp.Blogs, marshalBlog(blog))
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	})
}

func HandleBlogDelete(store *storage.Storage) http.Handler {
	type response struct {
		Blog jsonBlog `json:"blog"`
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

		err = store.Blog().Delete(blog)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Blog: marshalBlog(blog),
		}

		code := http.StatusOK
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	})
}
