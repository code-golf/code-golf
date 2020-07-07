package middleware

import (
	"database/sql"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// DatabaseHandler adds the database handle to the context.
func DatabaseHandler(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, session.Set(r, "database", db))
		})
	}
}
