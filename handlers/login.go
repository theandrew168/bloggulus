package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/theandrew168/bloggulus/models"

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

		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionID, err := generateSessionID()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		expiry := time.Now().AddDate(0, 0, 7)
		session := &models.Session{
			SessionID: sessionID,
			AccountID: account.AccountID,
			Expiry:    expiry,
		}
		session, err = app.Session.Create(r.Context(), session)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// create session_id cookie
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			Domain:   "",  // will default to the server's base domain
			Expires:  time.Unix(expiry.Unix() + 1, 0),  // round up to nearest second
			MaxAge:   int(time.Until(expiry).Seconds() + 1),  // round up to nearest second
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}

		// TODO: find a better way to check if running with TLS?
		if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
			cookie.Secure = true
		}

		// add headers to set the cookie
		w.Header().Add("Set-Cookie", cookie.String())
		w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
		w.Header().Add("Vary", "Cookie")

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

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
