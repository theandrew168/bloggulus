package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/models"
)

type followedBlog struct {
	Blog     *models.Blog
	Followed bool
}

type blogsData struct {
	Authed bool
	Blogs  []followedBlog
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
	if !authed {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	followed, err := app.Blog.ReadFollowedForUser(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	unfollowed, err := app.Blog.ReadUnfollowedForUser(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := &blogsData{
		Authed: authed,
	}

	for _, blog := range(followed) {
		data.Blogs = append(data.Blogs, followedBlog{blog, true})
	}

	for _, blog := range(unfollowed) {
		data.Blogs = append(data.Blogs, followedBlog{blog, false})
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
