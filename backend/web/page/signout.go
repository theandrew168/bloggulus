package page

import (
	"errors"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleSignoutForm(store *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie(util.SessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Delete the existing session cookie.
		cookie := util.NewExpiredCookie(util.SessionCookieName)
		http.SetCookie(w, &cookie)

		// Lookup the session by it's client-side session ID.
		session, err := store.Session().ReadBySessionID(sessionID.Value)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Redirect(w, r, "/", http.StatusSeeOther)
			default:
				http.Error(w, err.Error(), 500)
			}
			return
		}

		// Delete the session from the database.
		err = store.Session().Delete(session)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Redirect(w, r, "/", http.StatusSeeOther)
			default:
				http.Error(w, err.Error(), 500)
			}
			return
		}

		// Redirect back to the index page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
