package middleware

import "net/http"

// AdminArea enforces that an admin golfer is logged in.
func AdminArea(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if golfer := Golfer(r); golfer != nil && golfer.Admin {
			next.ServeHTTP(w, r)
		} else {
			// TODO Serve custom 403 without import cycle.
			w.WriteHeader(http.StatusForbidden)
		}
	})
}
