package routes

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"regexp"
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
	value, _ := url.QueryUnescape(r.PathValue(key))
	return value
}

func redir(templateURL string) http.HandlerFunc {
	re := regexp.MustCompile("{[^{}]*}")
	return func(w http.ResponseWriter, r *http.Request) {
		url := re.ReplaceAllStringFunc(templateURL, func(s string) string {
			// slice to remove the surrounding {}
			return param(r, s[1:len(s)-1])
		})
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
