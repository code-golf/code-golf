package middleware

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var public = map[string]bool{}

func init() {
	if err := filepath.Walk("public", func(file string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			public[file] = true
		}

		return err
	}); err != nil {
		panic(err)
	}
}

// Public serves any requested file from public/
func Public(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if name := path.Join("public", r.URL.Path); public[name] {
			http.ServeFile(w, r, name)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
