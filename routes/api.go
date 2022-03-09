package routes

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

//go:embed api.yml
var yml []byte

// API serves GET /api
// Use text/plain to always render in browser unlike the YML content types.
func API(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(yml)
}

// APICheevos serves GET /api/cheevos
func APICheevos(w http.ResponseWriter, _ *http.Request) {
	if err := json.NewEncoder(w).Encode(config.CheevoList); err != nil {
		panic(err)
	}
}

// APICheevo serves GET /api/cheevos/{cheevo}
func APICheevo(w http.ResponseWriter, r *http.Request) {
	if cheevo, ok := config.CheevoByID[param(r, "cheevo")]; !ok {
		APINotFound(w, r)
	} else if err := json.NewEncoder(w).Encode(cheevo); err != nil {
		panic(err)
	}
}

// APILangs serves GET /api/langs
func APILangs(w http.ResponseWriter, _ *http.Request) {
	if err := json.NewEncoder(w).Encode(config.LangList); err != nil {
		panic(err)
	}
}

// APILang serves GET /api/langs/{lang}
func APILang(w http.ResponseWriter, r *http.Request) {
	if lang, ok := config.LangByID[param(r, "lang")]; !ok {
		APINotFound(w, r)
	} else if err := json.NewEncoder(w).Encode(lang); err != nil {
		panic(err)
	}
}

// APINotFound serves an API 404.
func APINotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("null"))
}

// APIPanic serves GET /api/panic
func APIPanic(_ http.ResponseWriter, _ *http.Request) { panic("") }

// APISuggestionsGolfers serves GET /api/suggestions/golfers
func APISuggestionsGolfers(w http.ResponseWriter, r *http.Request) {
	var json []byte

	if err := session.Database(r).QueryRow(
		`WITH golfers AS (
		    SELECT login
		      FROM users
		     WHERE strpos(login, $1) > 0 AND login != $2
		  ORDER BY login
		     LIMIT 10
		) SELECT COALESCE(json_agg(login), '[]') FROM golfers`,
		r.FormValue("q"),
		r.FormValue("ignore"),
	).Scan(&json); err != nil {
		panic(err)
	}

	if _, err := w.Write(json); err != nil {
		panic(err)
	}
}
