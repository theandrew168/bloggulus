package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/model"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/storage"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

type jsonAccount struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func marshalAccount(account *model.Account) jsonAccount {
	a := jsonAccount{
		ID:       account.ID(),
		Username: account.Username(),
	}
	return a
}

func HandleAccountCreate(store *storage.Storage) http.Handler {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Account jsonAccount `json:"account"`
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

		e.CheckField(req.Username != "", "must be provided", "username")
		e.CheckField(len(req.Username) <= 500, "must not be more than 500 bytes long", "username")

		e.CheckField(req.Password != "", "must be provided", "password")
		e.CheckField(len(req.Password) >= 8, "must be at least 8 bytes long", "password")
		e.CheckField(len(req.Password) <= 72, "must not be more than 72 bytes long", "password")

		if !e.Valid() {
			util.FailedValidationResponse(w, r, e)
			return
		}

		account, err := model.NewAccount(req.Username, req.Password)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		err = store.Account().Create(account)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				e.AddField("already exists", "username")
				util.FailedValidationResponse(w, r, e)
			default:
				util.ServerErrorResponse(w, r, err)
			}

			return
		}

		resp := response{
			Account: marshalAccount(account),
		}

		code := http.StatusCreated
		err = util.WriteJSON(w, code, resp, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
	})
}
