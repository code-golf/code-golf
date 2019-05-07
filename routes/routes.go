package routes

import (
	"database/sql"
	"math/rand"
	"net/http"

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

	Router.GET("/", home)
	Router.GET("/about", about)
	Router.GET("/assets/:asset", asset)
	Router.GET("/callback", callback)
	Router.GET("/favicon.ico", asset)
	Router.GET("/logout", logout)
	Router.GET("/random", random)
	Router.GET("/robots.txt", robots)
	Router.GET("/scores/:hole/:lang", scores)
	Router.GET("/scores/:hole/:lang/mini", scoresMini)
	Router.GET("/stats", stats)
	Router.GET("/users/:user", user)

	for _, h := range holes {
		Router.GET("/"+h.ID, hole)
	}

	Router.POST("/solution", solution)

	Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Render(w, r, http.StatusNotFound, "404", nil)
	})
}

func random(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, holes[rand.Intn(len(holes))].ID, 302)
}

func robots(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNoContent)
}
