package util

import (
	"net/http"
	"time"
)

var (
	SessionIDCookieName = "session_id"
)

// Create a session (not permanent) cookie that expires when the user's session ends.
func NewSessionCookie(name, value string) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",  // This path makes the cookie apply to the whole site.
		Domain:   "",   // An empty domain will default to the server's base domain.
		Secure:   true, // Only send cookies on secure connections (includes localhost).
		HttpOnly: true, // Only send cookies via HTTP requests (not JS).
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}

// Create a permanent (not session) cookie with a given expiry.
func NewPermanentCookie(name, value string, expiry time.Time) http.Cookie {
	// Round the cookie's expiration up to nearest second.
	cookie := NewSessionCookie(name, value)
	cookie.Expires = time.Unix(expiry.Unix()+1, 0)
	cookie.MaxAge = int(time.Until(expiry).Seconds() + 1)
	return cookie
}

// Create a cookie that is instantly expired.
func NewExpiredCookie(name string) http.Cookie {
	cookie := NewSessionCookie(name, "")
	cookie.Expires = time.Unix(1, 0)
	cookie.MaxAge = -1
	return cookie
}
