package web

import (
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleRegisterPage() http.Handler {
	tmpl := page.NewRegister()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := page.RegisterData{}
		util.Render(w, r, http.StatusOK, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleRegisterForm(repo *repository.Repository) http.Handler {
	tmpl := page.NewRegister()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data.
		err := r.ParseForm()
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		// Pull out the expected form fields
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		// Validate the form values.
		v := util.NewValidator()
		v.CheckRequired("username", username)
		v.CheckRequired("password", password)

		// If the form isn't valid, re-render the template with existing input values.
		if !v.IsValid() {
			data := page.RegisterData{
				Username: username,
				Errors:   v,
			}
			util.Render(w, r, http.StatusBadRequest, func(w io.Writer) error {
				return tmpl.Render(w, data)
			})
			return
		}

		// Create a new account.
		account, err := model.NewAccount(username, password)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		// Save the new the account in the database.
		err = repo.Account().Create(account)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				// If a conflict occurs, re-render the form with an error.
				v.Add("username", "Username is already taken")
				data := page.RegisterData{
					Username: username,
					Errors:   v,
				}
				util.Render(w, r, http.StatusBadRequest, func(w io.Writer) error {
					return tmpl.Render(w, data)
				})
			default:
				util.InternalServerErrorResponse(w, r, err)
			}

			return
		}

		slog.Info("register",
			"account_id", account.ID(),
			"account_username", username,
		)

		session, sessionID, err := model.NewSession(account, util.SessionCookieTTL)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		err = repo.Session().Create(session)
		if err != nil {
			util.CreateErrorResponse(w, r, err)
			return
		}

		// Just set a session cookie (not a permanent one) after registration.
		cookie := util.NewSessionCookie(util.SessionCookieName, sessionID)
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
