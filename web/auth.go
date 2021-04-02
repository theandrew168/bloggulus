package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
)

var ErrNoSession = errors.New("user request doesn't have a valid session")

var SessionIDCookieName = "session_id"

func (app *Application) CheckSessionAccount(w http.ResponseWriter, r *http.Request) (int, error) {
	// check for session cookie
	sessionID, err := r.Cookie(SessionIDCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return 0, ErrNoSession
		} else {
			return 0, err
		}
	}

	// lookup session in the database
	session, err := app.Session.Read(r.Context(), sessionID.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// user has a session cookie but it's expired
			cookie := GenerateExpiredCookie(SessionIDCookieName)
			w.Header().Add("Set-Cookie", cookie.String())
			w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
			w.Header().Add("Vary", "Cookie")
			return 0, ErrNoSession
		} else {
			return 0, err
		}
	}

	// return ID of the account associated with the session
	return session.AccountID, nil
}

func GenerateSessionCookie(name, value string, expiry time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",  // applies to the whole site
		Domain:   "",  // will default to the server's base domain
		Expires:  time.Unix(expiry.Unix() + 1, 0),  // round up to nearest second
		MaxAge:   int(time.Until(expiry).Seconds() + 1),  // round up to nearest second
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}

func GenerateExpiredCookie(name string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",  // applies to the whole site
		Domain:   "",  // will default to the server's base domain
		Expires:  time.Unix(1, 0),  // expires now
		MaxAge:   -1,  // expires now
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}
