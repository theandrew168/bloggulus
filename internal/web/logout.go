package web

import (
	"log"
	"net/http"
)

func (app *Application) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// check for session cookie
	sessionID, err := r.Cookie(SessionIDCookieName)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// lookup session in the database
	session, err := app.Session.Read(r.Context(), sessionID.Value)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// delete session from the database
	err = app.Session.Delete(r.Context(), session.SessionID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// clear session cookie
	cookie := GenerateExpiredCookie(SessionIDCookieName)
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
