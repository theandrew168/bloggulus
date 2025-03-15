package web

import (
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

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
		account, isLoggedIn := util.GetContextAccount(r)
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
			BaseData: util.TemplateBaseData(r, w),
		}
		for _, blog := range blogs {
			data.Blogs = append(data.Blogs, page.BlogsBlogData{
				BaseData: util.TemplateBaseData(r, w),

				BlogForAccount: blog,
			})
		}

		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleBlogCreateForm(repo *repository.Repository, find *finder.Finder, syncService *service.SyncService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, isLoggedIn := util.GetContextAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		err := r.ParseForm()
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		feedURL := r.PostForm.Get("feedURL")

		// Check if the blog already exists.
		blog, err := repo.Blog().ReadByFeedURL(feedURL)
		if err == nil {
			// If it does, follow it for the current user.
			err = repo.AccountBlog().Create(account, blog)
			if err != nil {
				slog.Error("error following blog",
					"error", err.Error(),
					"account_id", account.ID(),
					"account_username", account.Username(),
					"blog_id", blog.ID(),
					"blog_title", blog.Title(),
				)
				return
			}

			slog.Info("blog followed",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"blog_id", blog.ID(),
				"blog_title", blog.Title(),
			)

			// Show a toast explaining that the blog already exists but is now being followed.
			cookie := util.NewSessionCookie(util.ToastCookieName, "This blog is now being followed!")
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		// At this point, the only "expected" error is ErrNotFound.
		if !errors.Is(err, postgres.ErrNotFound) {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		// Use the SyncService to add the new blog.
		// TODO: Make this respect graceful shutdowns.
		go func() {
			blog, err := syncService.SyncBlog(feedURL)
			if err != nil {
				slog.Error("error adding blog",
					"error", err.Error(),
					"feedURL", feedURL,
				)
				return
			}

			// TODO: Handle the case where the blog already exists (just follow it).

			slog.Info("blog added",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"blog_id", blog.ID(),
				"blog_title", blog.Title(),
			)

			err = repo.AccountBlog().Create(account, blog)
			if err != nil {
				slog.Error("error following blog",
					"error", err.Error(),
					"account_id", account.ID(),
					"account_username", account.Username(),
					"blog_id", blog.ID(),
					"blog_title", blog.Title(),
				)
				return
			}

			slog.Info("blog followed",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"blog_id", blog.ID(),
				"blog_title", blog.Title(),
			)
		}()

		// Show a toast explaining that the blog will be processed in the background.
		cookie := util.NewSessionCookie(util.ToastCookieName, "Once processed, this blog will be added and followed. Check back soon!")
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}

func HandleBlogFollowForm(repo *repository.Repository, find *finder.Finder) http.Handler {
	tmpl := page.NewBlogs()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, isLoggedIn := util.GetContextAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		err := r.ParseForm()
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			util.NotFoundResponse(w, r)
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
			data := page.BlogsBlogData{
				BaseData: util.TemplateBaseData(r, w),

				BlogForAccount: finder.BlogForAccount{
					ID:          blog.ID(),
					Title:       blog.Title(),
					SiteURL:     blog.SiteURL(),
					IsFollowing: true,
				},
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
		account, isLoggedIn := util.GetContextAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		err := r.ParseForm()
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		blogID, err := uuid.Parse(r.PathValue("blogID"))
		if err != nil {
			util.NotFoundResponse(w, r)
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
			data := page.BlogsBlogData{
				BaseData: util.TemplateBaseData(r, w),

				BlogForAccount: finder.BlogForAccount{
					ID:          blog.ID(),
					Title:       blog.Title(),
					SiteURL:     blog.SiteURL(),
					IsFollowing: false,
				},
			}
			util.Render(w, r, 200, func(w io.Writer) error {
				return tmpl.RenderBlog(w, data)
			})
			return
		}

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}
