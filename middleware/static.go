package middleware

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var public = map[string]bool{}

func init() {
	if err := filepath.WalkDir("public", func(file string, info fs.DirEntry, err error) error {
		if info != nil && !info.IsDir() {
			public[file] = true
		}

		return err
	}); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

// Static serves static files from public.
func Static(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serves any requested file from public/. These aren't cached.
		// Served without a prefix, e.g. GET /robots.txt.
		if file := path.Join("public", r.URL.Path); public[file] {
			http.ServeFile(w, r, file)
			return
		}

		next.ServeHTTP(w, r)
	})
}
