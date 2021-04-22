package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/theandrew168/bloggulus/rss"
	"github.com/theandrew168/bloggulus/model"
)

type followedBlog struct {
	Blog     *model.Blog
	Followed bool
}

type blogsData struct {
	Authed bool
	Blogs  []followedBlog
}

func (app *Application) HandleBlogs(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		feedURL := r.PostFormValue("feed_url")
		if feedURL == "" {
			log.Println("empty feed URL")
			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		blog, err := rss.ReadBlog(feedURL)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		_, err = app.Blog.Create(r.Context(), blog)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
		return
	}

	// TODO: sort followed
	followed, err := app.Blog.ReadFollowedForUser(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// TODO: sort unfollowed
	unfollowed, err := app.Blog.ReadUnfollowedForUser(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := &blogsData{
		Authed: authed,
	}

	for _, blog := range followed {
		data.Blogs = append(data.Blogs, followedBlog{blog, true})
	}

	for _, blog := range unfollowed {
		data.Blogs = append(data.Blogs, followedBlog{blog, false})
	}

	ts, err := template.ParseFiles("templates/blogs.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
