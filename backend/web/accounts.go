package web

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/theandrew168/bloggulus/backend/repository"
	"github.com/theandrew168/bloggulus/backend/web/ui"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func HandleAccountList(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accounts, err := repo.Account().List(100, 0)
		if err != nil {
			util.InternalServerErrorResponse(w, r, err)
			return
		}

		page := ui.AccountsPage(ui.AccountsPageData{
			PageLayoutData: util.GetPageLayoutData(r, w),

			Accounts: accounts,
		})
		util.Render(w, r, 200, func(w io.Writer) error {
			return page.Render(w)
		})
	})
}

func HandleAccountDeleteForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID, err := uuid.Parse(r.PathValue("accountID"))
		if err != nil {
			util.NotFoundResponse(w, r)
			return
		}

		account, err := repo.Account().Read(accountID)
		if err != nil {
			util.ReadErrorResponse(w, r, err)
			return
		}

		// Prevent accidental deletion of admin accounts.
		if account.IsAdmin() {
			slog.Info("admin account deletion attempt",
				"account_id", account.ID(),
				"account_username", account.Username(),
			)
			util.BadRequestResponse(w, r)
			return
		}

		err = repo.Account().Delete(account)
		if err != nil {
			util.DeleteErrorResponse(w, r, err)
			return
		}

		slog.Info("account deleted",
			"account_id", account.ID(),
			"account_username", account.Username(),
		)

		// Redirect back to the accounts page.
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
	})
}
