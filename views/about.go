package views

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type aboutData struct {
	Authed bool
}

func (app *Application) HandleAbout(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/about.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	authed := false

	// check for valid session
	sessionID, err := r.Cookie("session_id")
	if err == nil {
		// user does have a session_id cookie

		_, err := app.Session.Read(r.Context(), sessionID.Value)
		if err == nil {
			authed = true
		} else {
			// must be expired!
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
		}
	}

	data := &aboutData{
		Authed: authed,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
