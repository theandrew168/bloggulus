package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/theandrew168/bloggulus/models"
)

type blogsData struct {
	Authed bool
	Blogs  []*models.Blog
}

func (app *Application) HandleBlogs(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/blogs.html.tmpl", "templates/base.html.tmpl")
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

	var blogs []*models.Blog
	if authed {
		blogs, err = app.Blog.ReadAllForUser(r.Context(), accountID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		blogs, err = app.Blog.ReadAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	data := &blogsData{
		Authed: authed,
		Blogs:  blogs,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
