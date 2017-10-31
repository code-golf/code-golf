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
	Router.GET("/scores/:hole/:lang", middleware.Gzip(scores))
	Router.GET("/users/:user", middleware.Gzip(user))

	// TODO Create this with a loop.
	Router.GET("/99-bottles-of-beer", middleware.Gzip(hole))
	Router.GET("/arabic-to-roman-numerals", middleware.Gzip(hole))
	Router.GET("/e", middleware.Gzip(hole))
	Router.GET("/emirp-numbers", middleware.Gzip(hole))
	Router.GET("/evil-numbers", middleware.Gzip(hole))
	Router.GET("/fibonacci", middleware.Gzip(hole))
	Router.GET("/fizz-buzz", middleware.Gzip(hole))
	Router.GET("/happy-numbers", middleware.Gzip(hole))
	Router.GET("/odious-numbers", middleware.Gzip(hole))
	Router.GET("/pascals-triangle", middleware.Gzip(hole))
	Router.GET("/pernicious-numbers", middleware.Gzip(hole))
	Router.GET("/prime-numbers", middleware.Gzip(hole))
	Router.GET("/seven-segment", middleware.Gzip(hole))
	Router.GET("/sierpiński-triangle", middleware.Gzip(hole))
	Router.GET("/spelling-numbers", middleware.Gzip(hole))
	Router.GET("/π", middleware.Gzip(hole))

	Router.POST("/solution", middleware.Gzip(solution))

	Router.NotFound = http.HandlerFunc(print404)
}
