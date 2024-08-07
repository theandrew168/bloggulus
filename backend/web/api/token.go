package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

// when newly-created (and only then), tokens will include their plaintext value
type jsonNewToken struct {
	ID        uuid.UUID `json:"id"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func marshalNewToken(token *model.Token, value string) jsonNewToken {
	a := jsonNewToken{
		ID:        token.ID(),
		Value:     value,
		ExpiresAt: token.ExpiresAt(),
	}
	return a
}

// TODO: For when listing / reading tokens
// type jsonToken struct {
// 	ID        uuid.UUID `json:"id"`
// 	ExpiresAt time.Time `json:"expires_at"`
// }

// func marshalToken(token *model.Token) jsonToken {
// 	a := jsonToken{
// 		ID:        token.ID(),
// 		ExpiresAt: token.ExpiresAt(),
// 	}
// 	return a
// }

func HandleTokenCreate(store *storage.Storage) http.Handler {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Token jsonNewToken `json:"token"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := util.NewErrors()
		body := util.ReadBody(w, r)

		var req request
		err := util.ReadJSON(body, &req, true)
		if err != nil {
			util.BadRequestResponse(w, r)
			return
		}

		e.CheckField(req.Username != "", "Username must be provided", "username")
		e.CheckField(req.Password != "", "Password must be provided", "password")

		if !e.Valid() {
			util.FailedValidationResponse(w, r, e)
			return
		}

		// check that an account exists with the given username
		account, err := store.Account().ReadByUsername(req.Username)
		if err != nil {
			switch err {
			case postgres.ErrNotFound:
				e.Add("Invalid username or password")
				util.FailedValidationResponse(w, r, e)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		// check that the given password matches the account
		if !account.PasswordMatches(req.Password) {
			e.Add("Invalid username or password")
			util.FailedValidationResponse(w, r, e)
			return
		}

		token, value, err := model.NewToken(account, 24*time.Hour)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		err = store.Token().Create(token)
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
	})
}
