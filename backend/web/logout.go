package web

import (
	"errors"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleLogoutForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a session ID. If there isn't one, just redirect back home.
		sessionID, err := r.Cookie(util.SessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Delete the existing session cookie.
		cookie := util.NewExpiredCookie(util.SessionCookieName)
		http.SetCookie(w, &cookie)

		// Lookup the session by its client-side session ID.
		session, err := repo.Session().ReadBySessionID(sessionID.Value)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Redirect(w, r, "/", http.StatusSeeOther)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		// Delete the session from the database.
		err = repo.Session().Delete(session)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Redirect(w, r, "/", http.StatusSeeOther)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		// Redirect back to the index page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
