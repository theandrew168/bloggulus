package middleware

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func loginRedirectURL(path string) *url.URL {
	url, err := url.Parse("/login")
	if err != nil {
		panic(err)
	}

	// Update the URL with a "next" query param for post-auth redirecting.
	q := url.Query()
	q.Set("next", path)
	url.RawQuery = q.Encode()

	return url
}

func Authenticate(repo *repository.Repository) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check for a sessionID cookie.
			sessionID, err := r.Cookie(util.SessionCookieName)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Lookup the account linked to the session.
			account, err := repo.Account().ReadBySessionID(sessionID.Value)
			if err != nil {
				// If the user has an invalid / expired session cookie, delete it.
				if errors.Is(err, postgres.ErrNotFound) {
					cookie := util.NewExpiredCookie(util.SessionCookieName)
					http.SetCookie(w, &cookie)

					next.ServeHTTP(w, r)
					return
				}

				util.InternalServerErrorResponse(w, r, err)
				return
			}

			// If it exists, attach the account to the request context.
			r = util.SetContextAccount(r, account)

			next.ServeHTTP(w, r)
		})
	}
}

func AccountRequired() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request context has no account, then the user is not logged in (redirect).
			_, ok := util.GetContextAccount(r)
			if !ok {
				url := loginRedirectURL(r.URL.Path)
				http.Redirect(w, r, url.String(), http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AdminRequired() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request context has no account, then the user is not logged in (redirect).
			account, ok := util.GetContextAccount(r)
			if !ok {
				url := loginRedirectURL(r.URL.Path)
				http.Redirect(w, r, url.String(), http.StatusSeeOther)
				return
			}

			// If the account exists but is not an admin account, show a 403 Forbidden page.
			if !account.IsAdmin() {
				util.ForbiddenResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
