package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/theandrew168/bloggulus/internal/model"
)

var ErrNoSession = errors.New("user request doesn't have a valid session")

var (
	SessionIDCookieName = "session_id"
	SuccessCookieName   = "success"
	ErrorCookieName     = "error"
)

func (app *Application) CheckAccount(w http.ResponseWriter, r *http.Request) (*model.Account, error) {
	// check for session cookie
	sessionID, err := r.Cookie(SessionIDCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, ErrNoSession
		} else {
			return nil, err
		}
	}

	// lookup session in the database
	session, err := app.Session.Read(r.Context(), sessionID.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// user has a session cookie but it's expired
			cookie := GenerateExpiredCookie(SessionIDCookieName)
			http.SetCookie(w, cookie)
			return nil, ErrNoSession
		} else {
			return nil, err
		}
	}

	account := &session.Account
	return account, nil
}

func GenerateSessionCookie(name, value string, expiry time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",                                   // applies to the whole site
		Domain:   "",                                    // will default to the server's base domain
		Expires:  time.Unix(expiry.Unix()+1, 0),         // round up to nearest second
		MaxAge:   int(time.Until(expiry).Seconds() + 1), // round up to nearest second
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
		Path:     "/",             // applies to the whole site
		Domain:   "",              // will default to the server's base domain
		Expires:  time.Unix(1, 0), // expires now
		MaxAge:   -1,              // expires now
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}
