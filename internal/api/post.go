package api

import (
	"context"
	"net/http"
)

func (app *Application) HandlePost(w http.ResponseWriter, r *http.Request) {
	posts, err := app.storage.PostReadRecent(context.Background(), 20, 0)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, err.Error(), 500)
	}

	err = writeJSON(w, 200, posts, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, err.Error(), 500)
	}
}
