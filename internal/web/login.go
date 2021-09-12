package web

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/bloggulus/internal/core"
)

type loginData struct {
	Authed  bool
	Success string
	Error   string
}

func (app *Application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if username == "" || password == "" {
			expiry := time.Now().Add(time.Hour * 12)
			cookie := GenerateSessionCookie(ErrorCookieName, "Empty username or password", expiry)
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		account, err := app.Account.ReadByUsername(r.Context(), username)
		if err != nil {
			expiry := time.Now().Add(time.Hour * 12)
			cookie := GenerateSessionCookie(ErrorCookieName, "Invalid username or password", expiry)
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
		if err != nil {
			expiry := time.Now().Add(time.Hour * 12)
			cookie := GenerateSessionCookie(ErrorCookieName, "Invalid username or password", expiry)
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionID, err := GenerateSessionID()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		expiry := time.Now().AddDate(0, 0, 7)
		session := core.Session{
			SessionID: sessionID,
			Expiry:    expiry,
			Account:   account,
		}
		err = app.Session.Create(r.Context(), &session)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		// create session cookie
		cookie := GenerateSessionCookie(SessionIDCookieName, sessionID, expiry)
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFS(app.TemplatesFS, "login.html.tmpl", "base.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = app.CheckAccount(w, r)
	if err != nil {
		if err != core.ErrNotExist {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	authed := err == nil

	data := &loginData{
		Authed: authed,
	}

	// check for success cookie
	successCookie, err := r.Cookie(SuccessCookieName)
	if err == nil {
		data.Success = successCookie.Value
		cookie := GenerateExpiredCookie(SuccessCookieName)
		http.SetCookie(w, &cookie)
	}

	// check for error cookie
	errorCookie, err := r.Cookie(ErrorCookieName)
	if err == nil {
		data.Error = errorCookie.Value
		cookie := GenerateExpiredCookie(ErrorCookieName)
		http.SetCookie(w, &cookie)
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
