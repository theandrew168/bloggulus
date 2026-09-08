package web

import (
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/command"
	"github.com/theandrew168/bloggulus/backend/postgres"
	webquery "github.com/theandrew168/bloggulus/backend/query/web"
	"github.com/theandrew168/bloggulus/backend/value"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleBlogList(qry *webquery.Query) http.Handler {
	tmpl := page.NewBlogs()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, isLoggedIn := util.GetContextAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		var blogs []webquery.Blog
		var err error

		if account.IsAdmin {
			blogs, err = qry.Blog().ListAll(account.ID)
		} else {
			blogs, err = qry.Blog().ListVisible(account.ID)
		}

		if err != nil {
			util.ListErrorResponse(w, r, err)
			return
		}

		data := page.BlogsData{
			BaseData: util.GetTemplateBaseData(r, w),
		}
		for _, blog := range blogs {
			data.Blogs = append(data.Blogs, page.BlogsBlogData{
				BaseData: util.GetTemplateBaseData(r, w),

				Blog: blog,
			})
		}

		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

// TODO: This handler is pretty large. Can it be split and simplified?
func HandleBlogCreateForm(cmd *command.Command, qry *webquery.Query) http.Handler {
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

		feedURL, err := value.NewURL(r.PostForm.Get("feedURL"))
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		// Check if the blog already exists.
		blog, err := qry.Blog().ReadDetailsByFeedURL(feedURL)
		if err == nil {
			// If it does, follow it for the current user.
			err = cmd.Account().FollowBlog(account.ID, blog.ID)
			if err != nil {
				if !errors.Is(err, postgres.ErrConflict) {
					slog.Error("error following blog",
						"error", err.Error(),
						"account_id", account.ID,
						"account_username", account.Username,
						"blog_id", blog.ID,
						"blog_title", blog.Title,
					)
					return
				}
			}

			slog.Info("blog followed",
				"account_id", account.ID,
				"account_username", account.Username,
				"blog_id", blog.ID,
				"blog_title", blog.Title,
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
		// TODO: Make this respect graceful shutdowns (River Queue?)
		go func() {
			err := cmd.Sync().SyncBlog(feedURL)
			if err != nil {
				slog.Error("error adding blog",
					"error", err.Error(),
					"feedURL", feedURL,
				)
				return
			}

			// The blog _should_ exist not if SyncBlog finished without errors.
			// That means any errors here are fatal.
			blog, err = qry.Blog().ReadDetailsByFeedURL(feedURL)
			if err != nil {
				slog.Error("error reading blog",
					"error", err.Error(),
					"feedURL", feedURL,
				)
				return
			}

			slog.Info("blog added",
				"account_id", account.ID,
				"account_username", account.Username,
				"blog_id", blog.ID,
				"blog_title", blog.Title,
			)

			err = cmd.Account().FollowBlog(account.ID, blog.ID)
			if err != nil {
				if !errors.Is(err, postgres.ErrConflict) {
					slog.Error("error following blog",
						"error", err.Error(),
						"account_id", account.ID,
						"account_username", account.Username,
						"blog_id", blog.ID,
						"blog_title", blog.Title,
					)
					return
				}
			}
		}()

		// Show a toast explaining that the blog will be processed in the background.
		cookie := util.NewSessionCookie(util.ToastCookieName, "Once processed, this blog will be added and followed. Check back soon!")
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	})
}

func HandleBlogFollowForm(cmd *command.Command, qry *webquery.Query) http.Handler {
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

		err = cmd.Account().FollowBlog(account.ID, blogID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				util.BadRequestResponse(w, r)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		blog, err := qry.Blog().ReadDetailsByID(blogID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		// If the request came in via HTMX, re-render the individual blog row.
		if util.IsHTMXRequest(r) {
			data := page.BlogsBlogData{
				BaseData: util.GetTemplateBaseData(r, w),

				ID:          blog.ID,
				Title:       blog.Title,
				SiteURL:     blog.SiteURL,
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

func HandleBlogUnfollowForm(cmd *command.Command, qry *webquery.Query) http.Handler {
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

		err = cmd.Account().UnfollowBlog(account.ID, blogID)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				util.BadRequestResponse(w, r)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		blog, err := qry.Blog().ReadDetailsByID(blogID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		// If the request came in via HTMX, re-render the individual row.
		if util.IsHTMXRequest(r) {
			data := page.BlogsBlogData{
				BaseData: util.GetTemplateBaseData(r, w),

				ID:          blog.ID,
				Title:       blog.Title,
				SiteURL:     blog.SiteURL,
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
