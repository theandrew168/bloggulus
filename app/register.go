package app

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/theandrew168/bloggulus/model"

	"golang.org/x/crypto/bcrypt"
)

type registerData struct {
	Authed  bool
	Success string
	Error   string
}

func (app *Application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		username := r.PostFormValue("username")
		password1 := r.PostFormValue("password1")
		password2 := r.PostFormValue("password2")
		if username == "" || password1 == "" || password2 == "" {
			expiry := time.Now().Add(time.Hour * 12)
			cookie := GenerateSessionCookie(ErrorCookieName, "Empty username or password", expiry)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if password1 != password2 {
			expiry := time.Now().Add(time.Hour * 12)
			cookie := GenerateSessionCookie(ErrorCookieName, "Passwords do not match", expiry)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		account := &model.Account{
			Username: username,
			Password: string(hash),
		}
		account, err = app.Account.Create(r.Context(), account)
		if err != nil {
			if err == model.ErrExist {
				expiry := time.Now().Add(time.Hour * 12)
				cookie := GenerateSessionCookie(ErrorCookieName, "Failed to create account", expiry)
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			} else {
				log.Println(err)
				http.Error(w, err.Error(), 500)
				return
			}
		}

		expiry := time.Now().Add(time.Hour * 12)
		cookie := GenerateSessionCookie(SuccessCookieName, "Account successfully created!", expiry)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles("templates/register.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = app.CheckAccount(w, r)
	if err != nil {
		if err != ErrNoSession {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	authed := err == nil

	data := &registerData{
		Authed: authed,
	}

	// check for success cookie
	cookie, err := r.Cookie(SuccessCookieName)
	if err == nil {
		data.Success = cookie.Value
		cookie = GenerateExpiredCookie(SuccessCookieName)
		http.SetCookie(w, cookie)
	}

	// check for error cookie
	cookie, err = r.Cookie(ErrorCookieName)
	if err == nil {
		data.Error = cookie.Value
		cookie = GenerateExpiredCookie(ErrorCookieName)
		http.SetCookie(w, cookie)
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
