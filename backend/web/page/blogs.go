package page

import (
	_ "embed"
	"errors"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

//go:embed blogs.html
var blogsHTML string

type BlogsPageData struct {
	Blogs []finder.BlogForAccount
}

func HandleBlogsPage(find *finder.Finder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("page").Parse(blogsHTML)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		account, ok := util.ContextGetAccount(r)
		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		blogs, err := find.ListBlogsForAccount(account)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := BlogsPageData{
			Blogs: blogs,
		}
		tmpl.Execute(w, data)
	})
}

// TODO: Split the three intents into separate funcs.
// TODO: Make a helper for checking for HTMX requests.
func HandleBlogsForm(repo *repository.Repository, find *finder.Finder, syncService *service.SyncService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("page").Parse(blogsHTML)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		account, isLoggedIn := util.ContextGetAccount(r)
		if !isLoggedIn {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		intent := r.PostForm.Get("intent")
		if intent != "add" && intent != "follow" && intent != "unfollow" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if intent == "add" {
			feedURL := r.PostForm.Get("feedURL")

			// Check if the blog already exists. If it does, return its details.
			_, err := repo.Blog().ReadByFeedURL(feedURL)
			if err == nil {
				return
			}

			// At this point, the only "expected" error is ErrNotFound.
			if !errors.Is(err, postgres.ErrNotFound) {
				http.Error(w, err.Error(), 500)
				return
			}

			// Use the SyncService to add the new blog.
			blog, err := syncService.SyncBlog(feedURL)
			if err != nil {
				switch {
				case errors.Is(err, fetch.ErrUnreachableFeed):
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				default:
					http.Error(w, err.Error(), 500)
				}

				return
			}

			slog.Info("blog added",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"blog_id", blog.ID(),
				"blog_title", blog.Title(),
			)

			// Follow the blog and check for ErrConflict (already following).
			err = repo.AccountBlog().Create(account, blog)
			if err != nil {
				switch {
				case errors.Is(err, postgres.ErrConflict):
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				default:
					http.Error(w, err.Error(), 500)
				}

				return
			}

			slog.Info("blog followed",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"blog_id", blog.ID(),
				"blog_title", blog.Title(),
			)

			// If the request came in via HTMX, re-render the list of blogs.
			if r.Header.Get("HX-Request") != "" {
				blogs, err := find.ListBlogsForAccount(account)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				data := BlogsPageData{
					Blogs: blogs,
				}
				tmpl.ExecuteTemplate(w, "blogs", data)
				return
			}

			// This intent is done now.
			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		blogID, err := uuid.Parse(r.PostForm.Get("blogID"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		blog, err := repo.Blog().Read(blogID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			default:
				http.Error(w, err.Error(), 500)
			}

			return
		}

		if intent == "follow" {
			// Follow the blog and check for ErrConflict (already following).
			err = repo.AccountBlog().Create(account, blog)
			if err != nil {
				switch {
				case errors.Is(err, postgres.ErrConflict):
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				default:
					http.Error(w, err.Error(), 500)
				}

				return
			}

			slog.Info("blog followed",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"blog_id", blog.ID(),
				"blog_title", blog.Title(),
			)
		} else {
			// Unfollow the blog and check for ErrNotFound (already not following).
			err = repo.AccountBlog().Delete(account, blog)
			if err != nil {
				switch {
				case errors.Is(err, postgres.ErrNotFound):
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				default:
					http.Error(w, err.Error(), 500)
				}

				return
			}

			slog.Info("blog unfollowed",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"blog_id", blog.ID(),
				"blog_title", blog.Title(),
			)
		}

		// If the request came in via HTMX, re-render the individual row.
		if r.Header.Get("HX-Request") != "" {
			data := finder.BlogForAccount{
				ID:          blog.ID(),
				Title:       blog.Title(),
				IsFollowing: intent == "follow",
			}
			tmpl.ExecuteTemplate(w, "blog", data)
			return
		}

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}
