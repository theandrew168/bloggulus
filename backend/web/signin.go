package web

import (
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

func HandleSigninPage() http.Handler {
	tmpl := page.NewSignin()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := page.SigninData{}
		util.Render(w, r, http.StatusOK, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleSigninForm(repo *repository.Repository) http.Handler {
	tmpl := page.NewSignin()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data.
		err := r.ParseForm()
		if err != nil {
			util.BadRequestResponse(w, r)
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
			util.Render(w, r, http.StatusBadRequest, func(w io.Writer) error {
				return tmpl.Render(w, data)
			})
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
				util.Render(w, r, http.StatusBadRequest, func(w io.Writer) error {
					return tmpl.Render(w, data)
				})
			default:
				util.InternalServerErrorResponse(w, r, err)
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
			util.Render(w, r, http.StatusBadRequest, func(w io.Writer) error {
				return tmpl.Render(w, data)
			})
			return
		}

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
