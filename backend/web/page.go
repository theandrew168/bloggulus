package web

import (
	"net/http"

	"github.com/theandrew168/bloggulus/backend/repository"
)

func HandlePageList(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Render the pages template. It'll look similar to the blogs page w/ an input at the top.
	})
}

func HandlePageCreateForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Instantly return and show a toast. In a background goro, fetch
		// the page, parse out the title, strip out HTML, and create the
		// database rows (page and account_page).
	})
}

func HandlePageDeleteForm(repo *repository.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Just delete the account_page entry for this account + page. This
		// is because pages _could_ be added by multiple accounts and we
		// wouldn't wanna delete them out from under other users. If necessarh,
		// a service could be written that "garbage collects" dead pages.
	})
}
