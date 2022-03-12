package routes

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func cookie(r *http.Request, name string) (value string) {
	if c, _ := r.Cookie(name); c != nil {
		value = c.Value
	}
	return
}

// Generate a random string for CSP & OAuth state.
//
// The generated value SHOULD be at least 128 bits long (before encoding), and
// SHOULD be generated via a cryptographically secure random number generator.
// https://w3c.github.io/webappsec-csp/#security-nonces
func nonce() string {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(nonce)
}

func param(r *http.Request, key string) string {
	value, _ := url.QueryUnescape(chi.URLParam(r, key))
	return value
}

func redir(url string) http.HandlerFunc {
	return http.RedirectHandler(url, http.StatusPermanentRedirect).ServeHTTP
}
