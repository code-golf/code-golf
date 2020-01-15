package routes

import (
	"database/sql"
	"html/template"
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

	Router.GET("/", index)
	Router.GET("/about", about)
	Router.GET("/assets/:asset", asset)
	Router.GET("/callback", callback)
	Router.GET("/favicon.ico", asset)
	Router.GET("/feeds/:feed", Feeds)
	Router.GET("/ideas", ideas)
	Router.GET("/log-out", logOut)
	Router.GET("/random", random)
	Router.GET("/recent", recent)
	Router.GET("/robots.txt", robots)
	Router.GET("/scores/:hole/:lang", scores)
	Router.GET("/scores/:hole/:lang/:page", scores)
	Router.GET("/stats", stats)
	Router.GET("/users/:user", user)

	for _, h := range holes {
		Router.GET("/"+h.ID, hole)
	}

	Router.POST("/solution", solution)

	Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Render(w, r, http.StatusNotFound, "404", "", nil)
	})
}

func about(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Render(w, r, http.StatusOK, "about", "About", template.HTML(versionTable))
}

func logOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Set-Cookie", "__Host-user=;MaxAge=0;Path=/;Secure")

	http.Redirect(w, r, "/", http.StatusFound)
}

func random(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, holes[rand.Intn(len(holes))].ID, http.StatusFound)
}

func robots(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusNoContent)
}
