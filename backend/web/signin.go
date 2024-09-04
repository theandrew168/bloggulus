package web

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleSigninPage() http.Handler {
	tmpl := page.NewSignin()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := page.SigninData{}
		err := tmpl.Render(w, data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
}

func HandleSigninForm(repo *repository.Repository) http.Handler {
	tmpl := page.NewSignin()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data.
		err := r.ParseForm()
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
			data := page.SigninData{
				Username: username,
				Errors:   e,
			}
			err = tmpl.Render(w, data)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			return
		}

		account, err := repo.Account().ReadByUsername(username)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				e.Add("username", "Invalid username or password")
				e.Add("password", "Invalid username or password")
				data := page.SigninData{
					Username: username,
					Errors:   e,
				}
				err = tmpl.Render(w, data)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			default:
				http.Error(w, err.Error(), 500)
			}
			return
		}

		ok := account.PasswordMatches(password)
		if !ok {
			e.Add("username", "Invalid username or password")
			e.Add("password", "Invalid username or password")
			data := page.SigninData{
				Username: username,
				Errors:   e,
			}
			err = tmpl.Render(w, data)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
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
