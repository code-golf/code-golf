package middleware

import (
	"database/sql"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// DatabaseHandler adds the database handle to the context.
func DatabaseHandler(db, dbBeta *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if session.Beta(r) {
				r = session.Set(r, "database", dbBeta)
			} else {
				r = session.Set(r, "database", db)
			}

			next.ServeHTTP(w, r)
		})
	}
}
