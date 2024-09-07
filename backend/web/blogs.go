package web

import (
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/finder"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/service"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleBlogList(find *finder.Finder) http.Handler {
	tmpl := page.NewBlogs()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, isLoggedIn := util.ContextGetAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		blogs, err := find.ListBlogsForAccount(account)
		if err != nil {
			util.ListErrorResponse(w, r, err)
			return
		}

		data := page.BlogsData{
			Account: account,
			Blogs:   blogs,
		}

		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleBlogCreateForm(repo *repository.Repository, find *finder.Finder, syncService *service.SyncService) http.Handler {
	tmpl := page.NewBlogs()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		account, isLoggedIn := util.ContextGetAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		feedURL := r.PostForm.Get("feedURL")

		// Check if the blog already exists. If it does, nothing further is required.
		_, err = repo.Blog().ReadByFeedURL(feedURL)
		if err == nil {
			return
		}

		// At this point, the only "expected" error is ErrNotFound.
		if !errors.Is(err, postgres.ErrNotFound) {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		// Use the SyncService to add the new blog.
		blog, err := syncService.SyncBlog(feedURL)
		if err != nil {
			switch {
			case errors.Is(err, fetch.ErrUnreachableFeed):
				util.BadRequestResponse(w, r)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}

			return
		}

		slog.Info("blog added",
			"account_id", account.ID(),
			"account_username", account.Username(),
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		err = repo.AccountBlog().Create(account, blog)
		if err != nil {
			util.CreateErrorResponse(w, r, err)
			return
		}

		slog.Info("blog followed",
			"account_id", account.ID(),
			"account_username", account.Username(),
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		// If the request came in via HTMX, re-render the list of blogs.
		if util.IsHTMXRequest(r) {
			blogs, err := find.ListBlogsForAccount(account)
			if err != nil {
				util.ListErrorResponse(w, r, err)
				return
			}
			data := page.BlogsData{
				Account: account,
				Blogs:   blogs,
			}
			util.Render(w, r, 200, func(w io.Writer) error {
				return tmpl.RenderBlogs(w, data)
			})
			return
		}

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}

func HandleBlogFollowForm(repo *repository.Repository, find *finder.Finder) http.Handler {
	tmpl := page.NewBlogs()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		account, isLoggedIn := util.ContextGetAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		blogID, err := uuid.Parse(r.PostForm.Get("blogID"))
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		blog, err := repo.Blog().Read(blogID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		// Follow the blog and check for ErrConflict (already following).
		err = repo.AccountBlog().Create(account, blog)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				util.BadRequestResponse(w, r)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		slog.Info("blog followed",
			"account_id", account.ID(),
			"account_username", account.Username(),
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		// If the request came in via HTMX, re-render the individual blog row.
		if util.IsHTMXRequest(r) {
			data := finder.BlogForAccount{
				ID:          blog.ID(),
				Title:       blog.Title(),
				SiteURL:     blog.SiteURL(),
				IsFollowing: true,
			}
			util.Render(w, r, 200, func(w io.Writer) error {
				return tmpl.RenderBlog(w, data)
			})
			return
		}

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}

func HandleBlogUnfollowForm(repo *repository.Repository, find *finder.Finder) http.Handler {
	tmpl := page.NewBlogs()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		account, isLoggedIn := util.ContextGetAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		blogID, err := uuid.Parse(r.PostForm.Get("blogID"))
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		blog, err := repo.Blog().Read(blogID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		// Unfollow the blog and check for ErrNotFound (already not following).
		err = repo.AccountBlog().Delete(account, blog)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.BadRequestResponse(w, r)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		slog.Info("blog unfollowed",
			"account_id", account.ID(),
			"account_username", account.Username(),
			"blog_id", blog.ID(),
			"blog_title", blog.Title(),
		)

		// If the request came in via HTMX, re-render the individual row.
		if util.IsHTMXRequest(r) {
			data := finder.BlogForAccount{
				ID:          blog.ID(),
				Title:       blog.Title(),
				SiteURL:     blog.SiteURL(),
				IsFollowing: false,
			}
			util.Render(w, r, 200, func(w io.Writer) error {
				return tmpl.RenderBlog(w, data)
			})
			return
		}

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}
