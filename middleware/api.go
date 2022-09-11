package middleware

import "net/http"

// API adds CORS headers.
func API(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
