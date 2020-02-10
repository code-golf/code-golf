package routes

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

func db(r *http.Request) *sql.DB {
	return r.Context().Value("db").(*sql.DB)
}

func param(r *http.Request, key string) string {
	value, _ := url.QueryUnescape(chi.URLParam(r, key))
	return value
}
