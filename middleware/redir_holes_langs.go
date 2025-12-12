package middleware

import (
	"net/http"
	"strings"

	"github.com/code-golf/code-golf/config"
)

// RedirHolesLangs redirects old values for {hole} and {lang}.
func RedirHolesLangs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		holeID := r.PathValue("hole")
		langID := r.PathValue("lang")

		// Permanent redirects
		if redirect, ok := config.HoleRedirects[holeID]; ok {
			redir(w, r, holeID, redirect, http.StatusPermanentRedirect)
			return
		}

		if redirect, ok := config.LangRedirects[langID]; ok {
			redir(w, r, langID, redirect, http.StatusPermanentRedirect)
			return
		}

		// Aliases
		if alias, ok := config.HoleAliases[holeID]; ok {
			redir(w, r, holeID, alias, http.StatusTemporaryRedirect)
			return
		}

		// No redirect
		next.ServeHTTP(w, r)
	})
}

func redir(w http.ResponseWriter, r *http.Request, old, new string, code int) {
	// FIXME Consider using chi.RouteContext(r.Context()).RoutePattern()
	//       and filling in hole/lang, more robust than strings.Replace().
	path := strings.Replace(r.URL.Path, "/"+old, "/"+new, 1)

	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, path, code)
}
