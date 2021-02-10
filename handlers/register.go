package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func (app *Application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("attempted register from user: %s\n", r.PostFormValue("username"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
