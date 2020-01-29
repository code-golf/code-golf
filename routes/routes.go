package routes

import (
	"database/sql"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB
var Router = chi.NewRouter()

func init() {
	var err error

	if db, err = sql.Open("postgres", ""); err != nil {
		panic(err)
	}

	Router.Use(middleware.RedirectSlashes)

	Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		Render(w, r, http.StatusNotFound, "404", "", nil)
	})

	Router.Get("/", index)
	Router.Get("/{hole}", hole)
	Router.Get("/about", about)
	Router.Get("/assets/{asset}", asset)
	Router.Get("/callback", callback)
	Router.Get("/favicon.ico", asset)
	Router.Get("/feeds/{feed}", Feeds)
	Router.Get("/ideas", ideas)
	Router.Get("/log-out", logOut)
	Router.Get("/random", random)
	Router.Get("/recent", recent)
	Router.Get("/robots.txt", robots)
	Router.Get("/scores/{hole}/{lang}", scores)
	Router.Get("/scores/{hole}/{lang}/{suffix}", scores)
	Router.Post("/solution", solution)
	Router.Get("/stats", stats)
	Router.Get("/users/{user}", user)
}

func about(w http.ResponseWriter, r *http.Request) {
	Render(w, r, http.StatusOK, "about", "About", template.HTML(versionTable))
}

func logOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Set-Cookie", "__Host-user=;MaxAge=0;Path=/;Secure")

	http.Redirect(w, r, "/", http.StatusFound)
}

func random(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, holes[rand.Intn(len(holes))].ID, http.StatusFound)
}

func robots(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
