package page

import (
	_ "embed"
	"errors"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

//go:embed register.html
var registerHTML string

type RegisterPageData struct {
	Search   string
	Username string
	Errors   map[string]string
}

func HandleRegisterPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("register").Parse(registerHTML)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := RegisterPageData{}
		tmpl.Execute(w, data)
	})
}

func HandleRegisterForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("register").Parse(registerHTML)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Parse the form data.
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Pull out the expected form fields
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		// Validate the form values.
		e := util.NewErrors()
		e.CheckRequired("username", username)
		e.CheckRequired("password", password)

		// If the form isn't valid, re-render the template with existing input values.
		if !e.OK() {
			data := RegisterPageData{
				Username: username,
				Errors:   e,
			}
			tmpl.Execute(w, data)
			return
		}

		// Create a new account.
		account, err := model.NewAccount(username, password)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Save the new the account in the database.
		err = repo.Account().Create(account)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				// If a conflict occurs, re-render the form with an error.
				e.Add("username", "Username is already taken")
				data := RegisterPageData{
					Username: username,
					Errors:   e,
				}
				tmpl.Execute(w, data)
			default:
				http.Error(w, err.Error(), 500)
			}

			return
		}

		slog.Info("register", "username", username, "account_id", account.ID())

		session, sessionID, err := model.NewSession(account, util.SessionCookieTTL)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = repo.Session().Create(session)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Just set a session cookie (not a permanent one) after registration.
		cookie := util.NewSessionCookie(util.SessionCookieName, sessionID)
		http.SetCookie(w, &cookie)

		slog.Info("signin", "username", username, "account_id", account.ID(), "session_id", session.ID())

		// Redirect back to the index page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
