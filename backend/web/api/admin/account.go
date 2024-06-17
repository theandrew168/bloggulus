package admin

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/domain/admin"
	"github.com/theandrew168/bloggulus/backend/postgres"
	"github.com/theandrew168/bloggulus/backend/web/util"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

type jsonAccount struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func marshalAccount(account *admin.Account) jsonAccount {
	a := jsonAccount{
		ID:       account.ID(),
		Username: account.Username(),
	}
	return a
}

func (app *Application) handleAccountCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Account jsonAccount `json:"account"`
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
		v.Check(len(req.Username) <= 500, "username", "must not be more than 500 bytes long")

		v.Check(req.Password != "", "password", "must be provided")
		v.Check(len(req.Password) >= 8, "password", "must be at least 8 bytes long")
		v.Check(len(req.Password) <= 72, "password", "must not be more than 72 bytes long")

		if !v.Valid() {
			util.FailedValidationResponse(w, r, v.Errors())
			return
		}

		account, err := admin.NewAccount(req.Username, req.Password)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}

		err = app.store.Admin().Account().Create(account)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrConflict):
				v.AddError("username", "already exists")
				util.FailedValidationResponse(w, r, v.Errors())
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
	}
}
