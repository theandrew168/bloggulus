package app

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/bloggulus/model"
)

type registerData struct {
	Authed bool
}

func (app *Application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
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

		account := &model.Account{
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
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = app.CheckSessionAccount(w, r)
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

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
