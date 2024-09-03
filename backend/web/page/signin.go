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

//go:embed signin.html
var SigninHTML string

type SigninData struct {
	Username string
	Errors   map[string]string
}

func HandleSigninPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("page").Parse(SigninHTML)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := SigninData{}
		tmpl.Execute(w, data)
	})
}

func HandleSigninForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("page").Parse(SigninHTML)
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
			data := SigninData{
				Username: username,
				Errors:   e,
			}
			tmpl.Execute(w, data)
			return
		}

		account, err := repo.Account().ReadByUsername(username)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				e.Add("username", "Invalid username or password")
				e.Add("password", "Invalid username or password")
				data := SigninData{
					Username: username,
					Errors:   e,
				}
				tmpl.Execute(w, data)
			default:
				http.Error(w, err.Error(), 500)
			}
			return
		}

		ok := account.PasswordMatches(password)
		if !ok {
			e.Add("username", "Invalid username or password")
			e.Add("password", "Invalid username or password")
			data := SigninData{
				Username: username,
				Errors:   e,
			}
			tmpl.Execute(w, data)
			return
		}

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

		// Set a permanent cookie after signin.
		cookie := util.NewPermanentCookie(util.SessionCookieName, sessionID, util.SessionCookieTTL)
		http.SetCookie(w, &cookie)

		slog.Info("signin",
			"account_id", account.ID(),
			"account_username", username,
			"session_id", session.ID(),
		)

		// Redirect back to the index page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
