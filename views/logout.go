package views

import (
	"log"
	"net/http"
	"time"
)

func (app *Application) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// check for valid session
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, err := app.Session.Read(r.Context(), sessionID.Value)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = app.Session.Delete(r.Context(), session.SessionID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// delete existing session_id cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Domain:   "",  // will default to the server's base domain
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		Secure:   true,  // prod only
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	w.Header().Add("Set-Cookie", cookie.String())
	w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
	w.Header().Add("Vary", "Cookie")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
