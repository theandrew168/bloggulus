package web

import (
	"log"
	"net/http"
	"strconv"
)

func (app *Application) HandleFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
		return
	}

	// parse form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
		return
	}

	// pull out blog ID and convert to int
	blogID, err := strconv.Atoi(r.PostFormValue("blog_id"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
		return
	}

	// check for session cookie
	sessionID, err := r.Cookie(SessionIDCookieName)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
		return
	}

	// lookup session in the database
	session, err := app.Session.Read(r.Context(), sessionID.Value)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
		return
	}

	accountID := session.Account.AccountID
	log.Printf("account %d follow blog %d\n", accountID, blogID)
	app.AccountBlog.Follow(r.Context(), accountID, blogID)

	http.Redirect(w, r, "/blogs", http.StatusSeeOther)
	return
}
