package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/models"
)

type blogPost struct {
	Blog *models.Blog
	Post *models.Post
}

type indexData struct {
	Authed    bool
	BlogPosts []*blogPost
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html.tmpl", "templates/base.html.tmpl")
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

	var blogPosts []*blogPost
	if authed {
		posts, err := app.Post.ReadRecentForUser(r.Context(), accountID, 10)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, post := range posts {
			blog, err := app.Blog.Read(r.Context(), post.BlogID)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			blogPosts = append(blogPosts, &blogPost{blog, post})
		}
	} else {
		posts, err := app.Post.ReadRecent(r.Context(), 10)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, post := range posts {
			blog, err := app.Blog.Read(r.Context(), post.BlogID)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			blogPosts = append(blogPosts, &blogPost{blog, post})
		}
	}

	data := &indexData{
		Authed:    authed,
		BlogPosts: blogPosts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
