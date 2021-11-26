package api

import (
	"context"
	"net/http"
)

func (app *Application) HandleBlog(w http.ResponseWriter, r *http.Request) {
	blogs, err := app.storage.BlogReadAll(context.Background())
	if err != nil {
		app.logger.Println(err)
		http.Error(w, err.Error(), 500)
	}

	err = writeJSON(w, 200, blogs, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, err.Error(), 500)
	}
}
