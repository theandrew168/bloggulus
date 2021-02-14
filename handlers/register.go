package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/theandrew168/bloggulus/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterData struct {
	Authed bool
}

func (app *Application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		username := r.PostFormValue("username")
		password1 := r.PostFormValue("password1")
		password2 := r.PostFormValue("password2")
		if username == "" || password1 == "" || password2 == "" {
			log.Println("empty username or passwords")
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if password1 != password2 {
			log.Println("passwords don't match")
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		account := &models.Account{
			Username: username,
			Password: string(hash),
		}
		account, err = app.Account.Create(r.Context(), account)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles("templates/register.html.tmpl", "templates/base.html.tmpl")
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

	data := &RegisterData{
		Authed: authed,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
