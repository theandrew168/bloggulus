package handlers

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (app *Application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if username == "" || password == "" {
			log.Println("empty username or password")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		account, err := app.Account.ReadByUsername(r.Context(), username)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		log.Printf("attempted login from user: %v\n", account)
		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		log.Println("successful login!")
		log.Println("TODO: create session")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles("templates/login.html.tmpl")
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
