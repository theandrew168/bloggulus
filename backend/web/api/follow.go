package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleBlogFollow(store *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := util.NewErrors()

		account, ok := util.ContextGetAccount(r)
		if !ok {
			util.UnauthorizedResponse(w, r)
			return
		}

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

		err = store.AccountBlog().Create(account, blog)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				e.Add("You are already following this blog")
				util.FailedValidationResponse(w, r, e)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleBlogUnfollow(store *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := util.NewErrors()

		account, ok := util.ContextGetAccount(r)
		if !ok {
			util.UnauthorizedResponse(w, r)
			return
		}

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

		err = store.AccountBlog().Delete(account, blog)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				e.Add("You are not following this blog")
				util.FailedValidationResponse(w, r, e)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleBlogFollowing(store *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, ok := util.ContextGetAccount(r)
		if !ok {
			util.UnauthorizedResponse(w, r)
			return
		}

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

		count, err := store.AccountBlog().Count(account, blog)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		if count == 0 {
			util.NotFoundResponse(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
