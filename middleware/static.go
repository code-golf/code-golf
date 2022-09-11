package middleware

import (
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	dist   = ls("dist")
	public = ls("public")
)

func init() {
	if err := mime.AddExtensionType(".map", "application/json"); err != nil {
		panic(err)
	}
}

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

// Static serves static files from either dist or public.
func Static(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serves any requested file from dist/. These are immutable & cached.
		// Served with the dist prefix, e.g. GET /dist/js/foo-HASH.js.
		if file := r.URL.Path[1:]; dist[file] {
			w.Header().Set("Cache-Control", "max-age=31536000, public, immutable")
			w.Header().Set("Vary", "Accept-Encoding")

			// Identify with original filename, not pre-compressed filename.
			if ctype := mime.TypeByExtension(filepath.Ext(file)); ctype != "" {
				w.Header().Set("Content-Type", ctype)
			}

			// TODO Content negotiation github.com/golang/go/issues/19307.
			accept := r.Header.Get("Accept-Encoding")
			if strings.Contains(accept, "br") && dist[file+".br"] {
				w.Header().Set("Content-Encoding", "br")
				file += ".br"
			} else if strings.Contains(accept, "gzip") && dist[file+".gz"] {
				w.Header().Set("Content-Encoding", "gzip")
				file += ".gz"
			}

			http.ServeFile(w, r, file)
			return
		}

		// Serves any requested file from public/. These aren't cached.
		// Served without a prefix, e.g. GET /robots.txt.
		if file := path.Join("public", r.URL.Path); public[file] {
			http.ServeFile(w, r, file)
			return
		}

		next.ServeHTTP(w, r)
	})
}
