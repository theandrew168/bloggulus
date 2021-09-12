package web

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/theandrew168/bloggulus/internal/core"
)

var (
	SessionIDCookieName = "session_id"
	SuccessCookieName   = "success"
	ErrorCookieName     = "error"
)

func (app *Application) CheckAccount(w http.ResponseWriter, r *http.Request) (core.Account, error) {
	// check for session cookie
	sessionID, err := r.Cookie(SessionIDCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return core.Account{}, core.ErrNotExist
		} else {
			return core.Account{}, err
		}
	}

	// lookup session in the database
	session, err := app.Session.Read(r.Context(), sessionID.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// user has a session cookie but it's expired
			cookie := GenerateExpiredCookie(SessionIDCookieName)
			http.SetCookie(w, &cookie)
			return core.Account{}, core.ErrNotExist
		} else {
			return core.Account{}, err
		}
	}

	return session.Account, nil
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateSessionCookie(name, value string, expiry time.Time) http.Cookie {
	cookie := http.Cookie{
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

func GenerateExpiredCookie(name string) http.Cookie {
	cookie := http.Cookie{
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
