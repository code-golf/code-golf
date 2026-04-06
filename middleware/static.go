package middleware

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var public = ls("public")

func ls(dir string) map[string]bool {
	files := map[string]bool{}

	if err := filepath.WalkDir(dir, func(file string, info fs.DirEntry, err error) error {
		if info != nil && !info.IsDir() {
			files[file] = true
		}

		return err
	}); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	return files
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
