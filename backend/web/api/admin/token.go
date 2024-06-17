package admin

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/theandrew168/bloggulus/backend/model/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/web/util"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

// when newly-created (and only then), tokens will include their plaintext value
type jsonNewToken struct {
	ID        uuid.UUID `json:"id"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

func marshalNewToken(token *admin.Token, value string) jsonNewToken {
	a := jsonNewToken{
		ID:        token.ID(),
		Value:     value,
		ExpiresAt: token.ExpiresAt(),
	}
	return a
}

type jsonToken struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func marshalToken(token *admin.Token) jsonToken {
	a := jsonToken{
		ID:        token.ID(),
		ExpiresAt: token.ExpiresAt(),
	}
	return a
}

func (app *Application) handleTokenCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Token jsonNewToken `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		body := util.ReadBody(w, r)

		var req request
		err := util.ReadJSON(body, &req, true)
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		v.Check(req.Username != "", "username", "must be provided")
		v.Check(req.Password != "", "password", "must be provided")

		if !v.Valid() {
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		// check that an account exists with the given username
		account, err := app.store.Admin().Account().ReadByUsername(req.Username)
		if err != nil {
			switch err {
			case postgres.ErrNotFound:
				util.UnauthorizedResponse(w, r)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		// check that the given password matches the account
		if !account.PasswordMatches(req.Password) {
			util.UnauthorizedResponse(w, r)
			return
		}

		token, value, err := admin.NewToken(account, 24*time.Hour)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		err = app.store.Admin().Token().Create(token)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		resp := response{
			Token: marshalNewToken(token, value),
		}

		code := http.StatusCreated
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
}
