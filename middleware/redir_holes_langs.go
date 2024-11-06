package middleware

import (
	"net/http"
	"strings"

	"github.com/code-golf/code-golf/config"
)

// RedirHolesLangs redirects old values for {hole} and {lang}.
func RedirHolesLangs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newPath := r.URL.Path

		// FIXME Consider using chi.RouteContext(r.Context()).RoutePattern()
		//       and filling in hole/lang, more robust than strings.Replace().

		// Permanent redirects
		for _, hole := range config.AllHoleList {
			if newPath != r.URL.Path {
				break
			}
			for _, redirect := range hole.Redirects {
				if r.PathValue("hole") == redirect {
					newPath = strings.Replace(newPath, "/"+redirect, "/"+hole.ID, 1)
				}
			}
		}

		if r.PathValue("lang") == "perl6" {
			newPath = strings.Replace(newPath, "/perl6", "/raku", 1)
		}

		if newPath != r.URL.Path {
			if r.URL.RawQuery != "" {
				newPath += "?" + r.URL.RawQuery
			}

			http.Redirect(w, r, newPath, http.StatusPermanentRedirect)
			return
		}

		// Aliases
		for _, hole := range config.AllHoleList {
			if newPath != r.URL.Path {
				break
			}
			for _, alias := range hole.Aliases {
				if r.PathValue("hole") == alias {
					newPath = strings.Replace(newPath, "/"+alias, "/"+hole.ID, 1)
				}
			}
		}
		if newPath != r.URL.Path {
			if r.URL.RawQuery != "" {
				newPath += "?" + r.URL.RawQuery
			}

			http.Redirect(w, r, newPath, http.StatusTemporaryRedirect)
			return
		}

		// No redirect
		next.ServeHTTP(w, r)
	})
}
