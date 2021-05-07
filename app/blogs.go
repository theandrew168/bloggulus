package app

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/theandrew168/bloggulus/model"
	"github.com/theandrew168/bloggulus/rss"
	"github.com/theandrew168/bloggulus/task"
)

type followedBlog struct {
	Blog     *model.Blog
	Followed bool
}

type blogsData struct {
	Authed  bool
	Success string
	Error   string
	Blogs   []followedBlog
}

func (app *Application) HandleBlogs(w http.ResponseWriter, r *http.Request) {
	account, err := app.CheckAccount(w, r)
	if err != nil {
		if err != ErrNoSession {
			log.Println(err)
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
			http.Error(w, err.Error(), 500)
			return
		}

		feedURL := r.PostFormValue("feed_url")
		if feedURL == "" {
			expiry := time.Now().Add(time.Hour * 12)
			cookie := GenerateSessionCookie(ErrorCookieName, "Empty RSS / Atom feed URL", expiry)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		blog, err := rss.ReadBlog(feedURL)
		if err != nil {
			expiry := time.Now().Add(time.Hour * 12)
			cookie := GenerateSessionCookie(ErrorCookieName, "Invalid RSS / Atom feed", expiry)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/blogs", http.StatusSeeOther)
			return
		}

		blog, err = app.Blog.Create(r.Context(), blog)
		if err != nil {
			if err == model.ErrExist {
				// blog already exists!
				// look it up and link to the account
				blog, err = app.Blog.ReadByURL(r.Context(), feedURL)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), 500)
					return
				}

				err = app.AccountBlog.Follow(r.Context(), account.AccountID, blog.BlogID)
				if err != nil {
					if err != model.ErrExist {
						log.Println(err)
						http.Error(w, err.Error(), 500)
						return
					}
				}

				expiry := time.Now().Add(time.Hour * 12)
				cookie := GenerateSessionCookie(SuccessCookieName, "RSS / Atom feed already exists!", expiry)
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/blogs", http.StatusSeeOther)
				return
			} else {
				log.Println(err)
				http.Error(w, err.Error(), 500)
				return
			}
		}

		// link the blog to the account
		err = app.AccountBlog.Follow(r.Context(), account.AccountID, blog.BlogID)
		if err != nil {
			if err != model.ErrExist {
				log.Println(err)
				http.Error(w, err.Error(), 500)
				return
			}
		}

		// sync the new blog
		syncBlogs := task.SyncBlogs(app.Blog, app.Post)
		go syncBlogs.RunNow()

		expiry := time.Now().Add(time.Hour * 12)
		cookie := GenerateSessionCookie(SuccessCookieName, "RSS / Atom feed successfully added!", expiry)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/blogs", http.StatusSeeOther)
		return
	}

	// TODO: sort followed
	followed, err := app.Blog.ReadFollowedForUser(r.Context(), account.AccountID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// TODO: sort unfollowed
	unfollowed, err := app.Blog.ReadUnfollowedForUser(r.Context(), account.AccountID)
	if err != nil {
		log.Println(err)
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

	// check for success cookie
	cookie, err := r.Cookie(SuccessCookieName)
	if err == nil {
		data.Success = cookie.Value
		cookie = GenerateExpiredCookie(SuccessCookieName)
		http.SetCookie(w, cookie)
	}

	// check for error cookie
	cookie, err = r.Cookie(ErrorCookieName)
	if err == nil {
		data.Error = cookie.Value
		cookie = GenerateExpiredCookie(ErrorCookieName)
		http.SetCookie(w, cookie)
	}

	ts, err := template.ParseFiles("templates/blogs.html.tmpl", "templates/base.html.tmpl")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
