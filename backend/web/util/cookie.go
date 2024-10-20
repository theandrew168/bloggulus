package util

import (
	"net/http"
	"time"
)

var (
	// Permanents cookies.
	SessionCookieName = "bloggulus_session"
	SessionCookieTTL  = 7 * 24 * time.Hour

	// Session cookies.
	NextCookieName  = "bloggulus_next"
	StateCookieName = "bloggulus_state"
	ToastCookieName = "bloggulus_toast"
)

// Return a base cookie with secure defaults.
func NewBaseCookie() http.Cookie {
	cookie := http.Cookie{
		Path:     "/",  // This path makes the cookie apply to the whole site.
		Domain:   "",   // An empty domain will default to the server's base domain.
		Secure:   true, // Only send cookies on secure connections (includes localhost).
		HttpOnly: true, // Only send cookies via HTTP requests (not JS).
		// Don't send cookies with cross-site requests but include them when navigating
		// to the origin site from an external location (like when following a link).
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}

// Create a session (not permanent) cookie that expires when the user's session ends.
func NewSessionCookie(name, value string) http.Cookie {
	cookie := NewBaseCookie()
	cookie.Name = name
	cookie.Value = value
	return cookie
}

// Create a permanent (not session) cookie with a given time-to-live.
func NewPermanentCookie(name, value string, ttl time.Duration) http.Cookie {
	cookie := NewSessionCookie(name, value)
	cookie.MaxAge = int(ttl.Seconds())
	return cookie
}

// Create a cookie that is instantly expired.
func NewExpiredCookie(name string) http.Cookie {
	cookie := NewSessionCookie(name, "")
	cookie.MaxAge = -1
	return cookie
}
