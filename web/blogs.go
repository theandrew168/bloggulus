package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/models"
)

type blogsData struct {
	Authed bool
	Blogs  []*models.Blog
}

func (app *Application) HandleBlogs(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/blogs.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	accountID, err := app.CheckSessionAccount(w, r)
	if err != nil {
		if err != ErrNoSession {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	authed := err == nil

	var blogs []*models.Blog
	if authed {
		blogs, err = app.Blog.ReadAllForUser(r.Context(), accountID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		blogs, err = app.Blog.ReadAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	data := &blogsData{
		Authed: authed,
		Blogs:  blogs,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
