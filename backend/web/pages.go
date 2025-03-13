package web

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/fetch"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandlePageList(repo *repository.Repository) http.Handler {
	tmpl := page.NewPages()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, isLoggedIn := util.GetContextAccount(r)
		if !isLoggedIn {
			util.ForbiddenResponse(w, r)
			return
		}

		pages, err := repo.Page().ListByAccount(account, 100, 0)
		if err != nil {
			util.ListErrorResponse(w, r, err)
			return
		}

		data := page.PagesData{
			BaseData: util.TemplateBaseData(r, w),

			Pages: pages,
		}
		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

// Instantly return and show a toast. In a background goro, fetch
// the page, parse out the title, strip out HTML, and create the
// database rows (page and account_page).
func HandlePageCreateForm(repo *repository.Repository, pageFetcher fetch.PageFetcher) http.Handler {
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

		rawURL := r.PostForm.Get("url")

		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		if parsedURL.Scheme == "" {
			parsedURL.Scheme = "https"
		}

		pageURL := parsedURL.String()

		// If the page already exists, just follow it and return.
		page, err := repo.Page().ReadByURL(pageURL)
		if err == nil {
			// Follow the page and check for ErrConflict (already followed).
			err = repo.AccountPage().Create(account, page)
			if err != nil {
				if !errors.Is(err, postgres.ErrConflict) {
					util.InternalServerErrorResponse(w, r, err)
					return
				}
			}

			slog.Info("page followed",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"page_id", page.ID(),
				"page_url", page.URL(),
				"page_title", page.Title(),
			)

			// Show a toast explaining that the page will be processed in the background.
			cookie := util.NewSessionCookie(util.ToastCookieName, "This pageh as been added!")
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/pages", http.StatusSeeOther)
			return
		}

		if !errors.Is(err, postgres.ErrNotFound) {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		// Fetch the page and follow (if valid) in the background.
		go func() {
			request := fetch.FetchPageRequest{
				URL: pageURL,
			}
			response, err := pageFetcher.FetchPage(request)
			if err != nil {
				slog.Error("error fetching page",
					"error", err.Error(),
				)
				return
			}

			page, err = model.NewPage(pageURL, pageURL, response.Content)
			if err != nil {
				util.InternalServerErrorResponse(w, r, err)
				return
			}

			err = repo.Page().Create(page)
			if err != nil {
				util.InternalServerErrorResponse(w, r, err)
				return
			}

			slog.Info("page added",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"page_id", page.ID(),
				"page_url", page.URL(),
				"page_title", page.Title(),
			)

			// Follow the page and check for ErrConflict (already followed).
			err = repo.AccountPage().Create(account, page)
			if err != nil {
				if !errors.Is(err, postgres.ErrConflict) {
					util.InternalServerErrorResponse(w, r, err)
					return
				}
			}

			slog.Info("page followed",
				"account_id", account.ID(),
				"account_username", account.Username(),
				"page_id", page.ID(),
				"page_url", page.URL(),
				"page_title", page.Title(),
			)
		}()

		// Show a toast explaining that the page will be processed in the background.
		cookie := util.NewSessionCookie(util.ToastCookieName, "Once processed, this page will be added. Check back soon!")
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/pages", http.StatusSeeOther)
	})
}

// Just delete the account_page entry for this account + page. This
// is because pages _could_ be added by multiple accounts and we
// wouldn't wanna delete them out from under other users. If necessarh,
// a service could be written that "garbage collects" dead pages.
func HandlePageUnfollowForm(repo *repository.Repository) http.Handler {
	tmpl := page.NewPages()
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

		pageID, err := uuid.Parse(r.PathValue("pageID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		existingPage, err := repo.Page().Read(pageID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		// Unfollow the page and check for ErrNotFound (already not following).
		err = repo.AccountPage().Delete(account, existingPage)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				util.BadRequestResponse(w, r)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		slog.Info("page unfollowed",
			"account_id", account.ID(),
			"account_username", account.Username(),
			"page_id", existingPage.ID(),
			"page_url", existingPage.URL(),
			"page_title", existingPage.Title(),
		)

		// If the request came in via HTMX, re-render just the list of pages.
		if util.IsHTMXRequest(r) {
			pages, err := repo.Page().ListByAccount(account, 100, 0)
			if err != nil {
				util.ListErrorResponse(w, r, err)
				return
			}

			data := page.PagesData{
				BaseData: util.TemplateBaseData(r, w),

				Pages: pages,
			}
			util.Render(w, r, 200, func(w io.Writer) error {
				return tmpl.RenderPages(w, data)
			})
			return
		}

		http.Redirect(w, r, "/pages", http.StatusSeeOther)
	})
}
