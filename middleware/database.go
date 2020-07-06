package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

var databaseKey = key("database")

// DatabaseHandler adds the database handle to the context.
func DatabaseHandler(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), databaseKey, db)))
		})
	}
}

// Database gets the database handle from the request context.
func Database(r *http.Request) *sql.DB {
	return r.Context().Value(databaseKey).(*sql.DB)
}
