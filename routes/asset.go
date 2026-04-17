package routes

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/code-golf/code-golf/assets"
)

func init() {
	if err := mime.AddExtensionType(".map", "application/json"); err != nil {
		panic(err)
	}
}

// GET /assets/*
func assetGET(w http.ResponseWriter, r *http.Request) {
	file := param(r, "*")
	if !assets.Exists(file) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "max-age=31536000, public, immutable")
	w.Header().Set("Vary", "Accept-Encoding")

	// Identify with original filename, not pre-compressed filename.
	if ctype := mime.TypeByExtension(filepath.Ext(file)); ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	// TODO Content negotiation github.com/golang/go/issues/19307.
	accept := r.Header.Get("Accept-Encoding")
	if strings.Contains(accept, "br") && assets.Exists(file+".br") {
		w.Header().Set("Content-Encoding", "br")
		file += ".br"
	} else if strings.Contains(accept, "gzip") && assets.Exists(file+".gz") {
		w.Header().Set("Content-Encoding", "gzip")
		file += ".gz"
	}

	http.ServeFileFS(w, r, assets.Files, file)
}
