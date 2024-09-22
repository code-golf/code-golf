package middleware

import (
	"net/http"
	"strings"
)

// RedirHolesLangs redirects old values for {hole} and {lang}.
func RedirHolesLangs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newPath := r.URL.Path

		// FIXME Consider using chi.RouteContext(r.Context()).RoutePattern()
		//       and filling in hole/lang, more robust than strings.Replace().

		switch r.PathValue("hole") {
		case "billiard":
			newPath = strings.Replace(newPath, "/billiard", "/billiards", 1)
		case "eight-queens":
			newPath = strings.Replace(newPath, "/eight-queens", "/n-queens", 1)
		case "factorial-factorisation-ascii":
			newPath = strings.Replace(newPath, "/factorial-factorisation-ascii",
				"/factorial-factorisation", 1)
		case "grid-packing":
			newPath = strings.Replace(newPath, "/grid-packing", "/css-grid", 1)
		}

		if r.PathValue("lang") == "perl6" {
			newPath = strings.Replace(newPath, "/perl6", "/raku", 1)
		}

		if newPath == r.URL.Path {
			next.ServeHTTP(w, r)
		} else {
			if r.URL.RawQuery != "" {
				newPath += "?" + r.URL.RawQuery
			}

			http.Redirect(w, r, newPath, http.StatusPermanentRedirect)
		}
	})
}
