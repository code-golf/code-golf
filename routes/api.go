package routes

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"time"

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
	encodeJSON(w, config.CheevoList)
}

// GET /api/cheevos/{cheevo}
func apiCheevoGET(w http.ResponseWriter, r *http.Request) {
	encodeJSON(w, config.CheevoByID[param(r, "cheevo")])
}

// GET /api/golfers/{golfer}
func apiGolferGET(w http.ResponseWriter, r *http.Request) {
	golfer := &struct {
		Admin    bool      `json:"admin"`
		Country  *string   `json:"country"`
		ID       int       `json:"id"`
		Name     string    `json:"name"`
		Pronouns *string   `json:"pronouns"`
		Sponsor  bool      `json:"sponsor"`
		Started  time.Time `json:"started"`
	}{}

	if err := session.Database(r).Get(
		golfer,
		`SELECT admin, id, login name, pronouns, sponsor, started,
		        CASE WHEN show_country THEN country END country
		   FROM users
		  WHERE login = $1`,
		param(r, "golfer"),
	); errors.Is(err, sql.ErrNoRows) {
		golfer = nil
	} else if err != nil {
		panic(err)
	}

	encodeJSON(w, golfer)
}

// GET /api/holes
func apiHolesGET(w http.ResponseWriter, _ *http.Request) {
	encodeJSON(w, config.HoleList)
}

// GET /api/langs/{lang}
func apiHoleGET(w http.ResponseWriter, r *http.Request) {
	encodeJSON(w, config.HoleByID[param(r, "hole")])
}

// GET /api/langs
func apiLangsGET(w http.ResponseWriter, _ *http.Request) {
	encodeJSON(w, config.LangList)
}

// GET /api/langs/{lang}
func apiLangGET(w http.ResponseWriter, r *http.Request) {
	encodeJSON(w, config.LangByID[param(r, "lang")])
}

// GET /api/notes
func apiNotesGET(w http.ResponseWriter, r *http.Request) {
	notes := []struct {
		Hole string `json:"hole"`
		Lang string `json:"lang"`
		Note string `json:"note"`
	}{}

	if err := session.Database(r).Select(
		&notes,
		"SELECT hole, lang, note FROM notes WHERE user_id = $1",
		session.Golfer(r).ID,
	); err != nil {
		panic(err)
	}

	encodeJSON(w, notes)
}

// DELETE /api/notes/{hole}/{lang}
func apiNoteDELETE(w http.ResponseWriter, r *http.Request) {
	hole := config.HoleByID[param(r, "hole")]
	lang := config.LangByID[param(r, "lang")]
	if hole == nil || lang == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Database(r).MustExec(
		`DELETE FROM notes
		  WHERE user_id = $1 AND hole = $2 AND lang = $3`,
		session.Golfer(r).ID,
		hole.ID,
		lang.ID,
	)

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/notes/{hole}/{lang}
func apiNoteGET(w http.ResponseWriter, r *http.Request) {
	hole := config.HoleByID[param(r, "hole")]
	lang := config.LangByID[param(r, "lang")]
	if hole == nil || lang == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var note []byte

	if err := session.Database(r).Get(
		&note,
		`SELECT note
		   FROM notes
		  WHERE user_id = $1 AND hole = $2 AND lang = $3`,
		session.Golfer(r).ID,
		hole.ID,
		lang.ID,
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(note)
}

// PUT /api/notes/{hole}/{lang}
func apiNotePUT(w http.ResponseWriter, r *http.Request) {
	hole := config.HoleByID[param(r, "hole")]
	lang := config.LangByID[param(r, "lang")]
	note, _ := io.ReadAll(r.Body)

	if hole == nil || lang == nil || len(note) == 0 || len(note) >= 128*1024 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Only sponsors can create or update notes, the rest can still read & delete.
	golfer := session.Golfer(r)
	if !golfer.SponsorOrAdmin() {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	session.Database(r).MustExec(
		`INSERT INTO notes (user_id, hole, lang, note)
		      VALUES       (     $1,   $2,   $3,   $4)
		 ON CONFLICT       (user_id, hole, lang)
		   DO UPDATE SET note = excluded.note`,
		golfer.ID,
		hole.ID,
		lang.ID,
		note,
	)

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/panic
func apiPanicGET(_ http.ResponseWriter, _ *http.Request) { panic("") }

// GET /api/solutions-log
func apiSolutionsLogGET(w http.ResponseWriter, r *http.Request) {
	hole := config.HoleByID[r.FormValue("hole")]
	lang := config.LangByID[r.FormValue("lang")]
	if hole == nil || lang == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rows []struct {
		Bytes     int       `json:"bytes"`
		Chars     int       `json:"chars"`
		Golfer    string    `json:"golfer"`
		Hole      string    `json:"hole"`
		Lang      string    `json:"lang"`
		Scoring   string    `json:"scoring"`
		Submitted time.Time `json:"submitted"`
	}
	if err := session.Database(r).SelectContext(
		r.Context(),
		&rows,
		` SELECT bytes, chars, login golfer, hole, lang, scoring, submitted
		    FROM solutions_log
		    JOIN users ON id = user_id
		   WHERE hole = $1 AND lang = $2
		ORDER BY submitted, scoring`,
		hole.ID,
		lang.ID,
	); err != nil {
		panic(err)
	}

	encodeJSON(w, rows)
}

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

// GET /api/wiki/{slug}
func apiWikiPageGET(w http.ResponseWriter, r *http.Request) {
	page := &struct {
		HTML string `json:"content"`
		Name string `json:"title"`
	}{}

	if err := session.Database(r).Get(
		page,
		"SELECT html, name FROM wiki WHERE slug = $1",
		param(r, "*"),
	); errors.Is(err, sql.ErrNoRows) {
		page = nil
	} else if err != nil {
		panic(err)
	}

	encodeJSON(w, page)
}

func encodeJSON(w http.ResponseWriter, v any) {
	if v == nil || reflect.ValueOf(v).IsNil() {
		w.WriteHeader(http.StatusNotFound)
	} else if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}
