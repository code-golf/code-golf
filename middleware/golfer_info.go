package middleware

import (
	"net/http"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/session"
)

// GolferInfo adds the GolferInfo object handle to the context.
func GolferInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		info := golfer.GetInfo(session.Database(r), name)

		switch {
		case info == nil:
			w.WriteHeader(http.StatusNotFound)
		case info.Name != name:
			// TODO Handle /holes suffix.
			path := "/golfers/" + info.Name
			if r.URL.RawQuery != "" {
				path += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, path, http.StatusTemporaryRedirect)
		default:
			session.Get(r).GolferInfo = info
			next.ServeHTTP(w, r)
		}
	})
}
