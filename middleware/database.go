package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Database adds the database handle to the context.
func Database(db *pgxpool.Pool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, session.Set(r, "database", db))
		})
	}
}
