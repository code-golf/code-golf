package routes

import (
	"database/sql"
	"net/http"

	"github.com/jraspass/code-golf/middleware"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var db *sql.DB
var Router = httprouter.New()

func init() {
	var err error

	if db, err = sql.Open("postgres", ""); err != nil {
		panic(err)
	}

	Router.GET("/", middleware.Gzip(home))
	Router.GET("/about", middleware.Gzip(about))
	Router.GET("/assets/:asset", asset)
	Router.GET("/callback", callback)
	Router.GET("/favicon.ico", asset)
	Router.GET("/logout", middleware.Gzip(logout))
	Router.GET("/scores", middleware.Gzip(scores))
	Router.GET("/scores/*criteria", middleware.Gzip(scores))
	Router.GET("/users/:user", middleware.Gzip(user))

	for name := range preambles {
		Router.GET("/" + name, middleware.Gzip(hole))
	}

	Router.POST("/solution", middleware.Gzip(solution))

	Router.NotFound = http.HandlerFunc(print404)
}
