package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/theandrew168/bloggulus/internal/core"
)

func (app *Application) HandleFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Add("Allow", "POST")
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	// parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// pull out blog ID and convert to int
	blogID, err := strconv.Atoi(r.PostFormValue("blog_id"))
	if err != nil {
		log.Println(err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	// lookup account
	account, err := app.CheckAccount(w, r)
	if err != nil {
		if err == ErrNoSession {
			status := http.StatusUnauthorized
			http.Error(w, http.StatusText(status), status)
		} else {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	// link the blog to the account
	err = app.Follow.Follow(r.Context(), account.AccountID, blogID)
	if err != nil {
		if err != core.ErrExist {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}

	http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	return
}
