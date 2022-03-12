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

// GET /api, text/plain rather that YML content types to render in browser.
func apiGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(yml)
}

// GET /api/cheevos
func apiCheevosGET(w http.ResponseWriter, _ *http.Request) {
	if err := json.NewEncoder(w).Encode(config.CheevoList); err != nil {
		panic(err)
	}
}

// GET /api/cheevos/{cheevo}
func apiCheevoGET(w http.ResponseWriter, r *http.Request) {
	if cheevo, ok := config.CheevoByID[param(r, "cheevo")]; !ok {
		w.WriteHeader(http.StatusNotFound)
	} else if err := json.NewEncoder(w).Encode(cheevo); err != nil {
		panic(err)
	}
}

// GET /api/langs
func apiLangsGET(w http.ResponseWriter, _ *http.Request) {
	if err := json.NewEncoder(w).Encode(config.LangList); err != nil {
		panic(err)
	}
}

// GET /api/langs/{lang}
func apiLangGET(w http.ResponseWriter, r *http.Request) {
	if lang, ok := config.LangByID[param(r, "lang")]; !ok {
		w.WriteHeader(http.StatusNotFound)
	} else if err := json.NewEncoder(w).Encode(lang); err != nil {
		panic(err)
	}
}

// GET /api/panic
func apiPanicGET(_ http.ResponseWriter, _ *http.Request) { panic("") }

// GET /api/suggestions/golfers
func apiSuggestionsGolfersGET(w http.ResponseWriter, r *http.Request) {
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
