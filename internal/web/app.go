package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/theandrew168/bloggulus/internal/core"
)

type Application struct {
	Account     core.AccountStorage
	AccountBlog core.AccountBlogStorage
	Blog        core.BlogStorage
	Post        core.PostStorage
	Session     core.SessionStorage
}

func (app *Application) Router() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", app.HandleIndex)
	router.HandlerFunc("GET", "/login", app.HandleLogin)
	router.HandlerFunc("POST", "/login", app.HandleLogin)
	router.HandlerFunc("POST", "/logout", app.HandleLogout)
	router.HandlerFunc("GET", "/blogs", app.HandleBlogs)
	router.HandlerFunc("POST", "/blogs", app.HandleBlogs)
	router.HandlerFunc("POST", "/follow", app.HandleFollow)
	router.HandlerFunc("POST", "/unfollow", app.HandleUnfollow)
	router.HandlerFunc("GET", "/register", app.HandleRegister)
	router.HandlerFunc("POST", "/register", app.HandleRegister)
	router.ServeFiles("/static/*filepath", http.Dir("./static"))
	return router
}
