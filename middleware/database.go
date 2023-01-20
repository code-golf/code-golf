package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
	"github.com/jmoiron/sqlx"
)

// Database adds the database handle to the context.
func Database(db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, session.Set(r, "database", db))
		})
	}
}
