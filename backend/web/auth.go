package web

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/random"
	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"

	"golang.org/x/oauth2"
)

type FetchUserID func(client *http.Client) (string, error)

func FetchGithubUserID(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		slog.Error("failed to obtain user information", "error", err.Error())
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read user information", "error", err.Error())
		return "", err
	}

	// Combine the provider and ID to create a unique identifier across all
	// OAuth services (like "github_123456" or "google_123456"). Then, hash
	// that ID before using as the account's username.
	type userinfo struct {
		ID json.Number `json:"id"`
	}

	var user userinfo
	err = json.Unmarshal(body, &user)
	if err != nil {
		slog.Error("failed to parse user information", "error", err.Error())
		return "", err
	}

	userID := user.ID.String()
	if userID == "" {
		slog.Error("failed to obtain user information")
		return "", err
	}

	userID = "github_" + userID
	return userID, nil
}

func FetchGoogleUserID(client *http.Client) (string, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		slog.Error("failed to obtain user information", "error", err.Error())
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read user information", "error", err.Error())
		return "", err
	}

	// Combine the provider and ID to create a unique identifier across all
	// OAuth services (like "github_123456" or "google_123456"). Then, hash
	// that ID before using as the account's username.
	type userinfo struct {
		ID string `json:"id"`
	}

	var user userinfo
	err = json.Unmarshal(body, &user)
	if err != nil {
		slog.Error("failed to parse user information", "error", err.Error())
		return "", err
	}

	userID := user.ID
	if userID == "" {
		slog.Error("failed to obtain user information")
		return "", err
	}

	userID = "google_" + userID
	return userID, nil
}

func HandleSignIn(enableDebugAuth bool) http.Handler {
	tmpl := page.NewSignIn()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a "next" query param for post-auth redirecting.
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}

		// Store "next" URL in a session cookie.
		cookie := util.NewSessionCookie(util.NextCookieName, next)
		http.SetCookie(w, &cookie)

		data := page.SignInData{
			BaseData: util.TemplateBaseData(r, w),

			EnableDebugAuth: enableDebugAuth,
		}
		util.Render(w, r, http.StatusOK, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleOAuthSignIn(conf *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := random.BytesBase64(16)
		if err != nil {
			panic(err)
		}

		cookie := util.NewSessionCookie(util.StateCookieName, state)
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, conf.AuthCodeURL(state), http.StatusFound)
	})
}

func HandleOAuthCallback(
	secretKey string,
	repo *repository.Repository,
	conf *oauth2.Config,
	fetchUserID FetchUserID,
) http.Handler {
	// TODO: Replace the 400s with sign in page re-renders.
	// tmpl := page.NewSignIn()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clear out the state expiredStateCookie.
		expiredStateCookie := util.NewExpiredCookie(util.StateCookieName)
		http.SetCookie(w, &expiredStateCookie)

		state, err := r.Cookie(util.StateCookieName)
		if err != nil {
			slog.Error("state not found")
			util.BadRequestResponse(w, r)
			return
		}

		if r.URL.Query().Get("state") != state.Value {
			slog.Error("state did not match")
			util.BadRequestResponse(w, r)
			return
		}

		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			slog.Error("failed to exchange code for access token", "error", err.Error())
			util.BadRequestResponse(w, r)
			return
		}

		client := conf.Client(context.Background(), token)
		userID, err := fetchUserID(client)
		if err != nil {
			slog.Error("failed to fetch user ID", "error", err.Error())
			util.BadRequestResponse(w, r)
			return
		}

		username := util.HashUserID(userID, secretKey)
		account, err := repo.Account().ReadByUsername(username)
		if err != nil {
			if !errors.Is(err, postgres.ErrNotFound) {
				util.InternalServerErrorResponse(w, r, err)
				return
			}

			// We need to create a new account at this point.
			account, err = model.NewAccount(username)
			if err != nil {
				util.InternalServerErrorResponse(w, r, err)
				return
			}

			err = repo.Account().Create(account)
			if err != nil {
				slog.Error("failed create user account", "error", err.Error())
				util.BadRequestResponse(w, r)
				return
			}

			slog.Info("account created",
				"account_id", account.ID(),
			)
		}

		// Create a new session for the account.
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

		// Set a permanent cookie after sign in.
		sessionCookie := util.NewPermanentCookie(util.SessionCookieName, sessionID, util.SessionCookieTTL)
		http.SetCookie(w, &sessionCookie)

		slog.Info("account signed in",
			"account_id", account.ID(),
			"session_id", session.ID(),
		)

		next := "/"
		nextCookie, err := r.Cookie(util.NextCookieName)
		if err == nil {
			next = nextCookie.Value

			expiredNextCookie := util.NewExpiredCookie(util.NextCookieName)
			http.SetCookie(w, &expiredNextCookie)
		}

		http.Redirect(w, r, next, http.StatusFound)
	})
}

func HandleDebugSignIn(secretKey string, repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a random userID for the debug sign in.
		userID, err := random.BytesBase64(16)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		userID = "debug_" + userID
		username := util.HashUserID(userID, secretKey)

		account, err := repo.Account().ReadByUsername(username)
		if err != nil {
			if !errors.Is(err, postgres.ErrNotFound) {
				util.InternalServerErrorResponse(w, r, err)
				return
			}

			// We need to create a new account at this point.
			account, err = model.NewAccount(username)
			if err != nil {
				util.InternalServerErrorResponse(w, r, err)
				return
			}

			err = repo.Account().Create(account)
			if err != nil {
				slog.Error("failed create user account", "error", err.Error())
				util.BadRequestResponse(w, r)
				return
			}

			slog.Info("account created",
				"account_id", account.ID(),
			)
		}

		// Create a new session for the account.
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

		// Set a permanent cookie after sign in.
		sessionCookie := util.NewPermanentCookie(util.SessionCookieName, sessionID, util.SessionCookieTTL)
		http.SetCookie(w, &sessionCookie)

		slog.Info("account signed in",
			"account_id", account.ID(),
			"session_id", session.ID(),
		)

		next := "/"
		nextCookie, err := r.Cookie(util.NextCookieName)
		if err == nil {
			next = nextCookie.Value

			expiredNextCookie := util.NewExpiredCookie(util.NextCookieName)
			http.SetCookie(w, &expiredNextCookie)
		}

		http.Redirect(w, r, next, http.StatusFound)
	})
}

func HandleSignOutForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a session ID. If there isn't one, just redirect back home.
		sessionID, err := r.Cookie(util.SessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Delete the existing session cookie.
		cookie := util.NewExpiredCookie(util.SessionCookieName)
		http.SetCookie(w, &cookie)

		// Lookup the session by its client-side session ID.
		session, err := repo.Session().ReadBySessionID(sessionID.Value)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Redirect(w, r, "/", http.StatusSeeOther)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		// Delete the session from the database.
		err = repo.Session().Delete(session)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				http.Redirect(w, r, "/", http.StatusSeeOther)
			default:
				util.InternalServerErrorResponse(w, r, err)
			}
			return
		}

		// Redirect back to the index page.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}