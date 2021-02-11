package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/models"

	"golang.org/x/crypto/bcrypt"
)

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

	ts, err := template.ParseFiles("templates/register.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
