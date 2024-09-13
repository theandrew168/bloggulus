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

func HandleLoginPage() http.Handler {
	tmpl := page.NewLogin()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := page.LoginData{
			BaseData: util.TemplateBaseData(r, w),
		}
		util.Render(w, r, http.StatusOK, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleLoginForm(repo *repository.Repository) http.Handler {
	tmpl := page.NewLogin()
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
		v := util.NewValidator()
		v.CheckRequired("username", username)
		v.CheckRequired("password", password)

		// If the form isn't valid, re-render the template with existing input values.
		if !v.IsValid() {
			data := page.LoginData{
				Username: username,
				Errors:   v,
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
				v.Add("username", "Invalid username or password")
				v.Add("password", "Invalid username or password")
				data := page.LoginData{
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

		ok := account.PasswordMatches(password)
		if !ok {
			v.Add("username", "Invalid username or password")
			v.Add("password", "Invalid username or password")
			data := page.LoginData{
				Username: username,
				Errors:   v,
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

		// Set a permanent cookie after login.
		cookie := util.NewPermanentCookie(util.SessionCookieName, sessionID, util.SessionCookieTTL)
		http.SetCookie(w, &cookie)

		slog.Info("login",
			"account_id", account.ID(),
			"account_username", username,
			"session_id", session.ID(),
		)

		// Redirect back to the index page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
