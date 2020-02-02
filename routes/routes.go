package routes

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error

	if db, err = sql.Open("postgres", ""); err != nil {
		panic(err)
	}
}

func param(r *http.Request, key string) string {
	value, _ := url.QueryUnescape(chi.URLParam(r, key))
	return value
}
