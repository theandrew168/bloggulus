package web

import (
	"io"
	"net/http"
	"uuid"

	"github.com/theandrew168/bloggulus/backend/command"
	webquery "github.com/theandrew168/bloggulus/backend/query/web"
	"github.com/theandrew168/bloggulus/backend/web/page"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleAccountList(qry *webquery.Query) http.Handler {
	tmpl := page.NewAccounts()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accounts, err := qry.Account().List()
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		data := page.AccountsData{
			BaseData: util.GetTemplateBaseData(r, w),

			Accounts: accounts,
		}
		util.Render(w, r, 200, func(w io.Writer) error {
			return tmpl.Render(w, data)
		})
	})
}

func HandleAccountDeleteForm(cmd *command.Command) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID, err := uuid.Parse(r.PathValue("accountID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		err = cmd.Account().DeleteAccount(accountID)
		if err != nil {
			util.DeleteErrorResponse(w, r, err)
			return
		}

		// Redirect back to the accounts page.
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	})
}
