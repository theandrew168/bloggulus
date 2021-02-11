package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/theandrew168/bloggulus/models"
)

type IndexData struct {
	Authed bool
	Posts  []*models.SourcedPost
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	authed := false
	accountID := 0

	// check for valid session
	sessionID, err := r.Cookie("session_id")
	if err == nil {
		// user does have a session_id cookie

		session, err := app.Session.Read(r.Context(), sessionID.Value)
		if err == nil {
			authed = true
			accountID = session.AccountID
		} else {
			// must be expired!
			// delete existing session_id cookie
			cookie := http.Cookie{
				Name:     "session_id",
				Path:     "/",
				Domain:   "",  // will default to the server's base domain
				Expires:  time.Unix(1, 0),
		//		Secure:   true,  // prod only
				MaxAge:   -1,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			w.Header().Add("Set-Cookie", cookie.String())
			w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
			w.Header().Add("Vary", "Cookie")
		}
	}

	var posts []*models.SourcedPost
	if authed {
		posts, err = app.SourcedPost.ReadRecentForUser(r.Context(), accountID, 20)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		posts, err = app.SourcedPost.ReadRecent(r.Context(), 20)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	data := &IndexData{
		Authed: authed,
		Posts:  posts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
