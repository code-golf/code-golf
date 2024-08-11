package routes

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
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

// GET /mini-rankings/{hole}/{lang}/{scoring:bytes|chars}/{view:top|me|following}
func apiMiniRankingsGET(w http.ResponseWriter, r *http.Request) {
	limit := 7
	if r.FormValue("long") == "1" {
		limit = 99
	}

	var (
		holeID  = param(r, "hole")
		langID  = param(r, "lang")
		scoring = param(r, "scoring")
		view    = param(r, "view")
	)

	// No need to check scoring & view, they're covered by chi.
	hole := config.AllHoleByID[holeID]
	lang := config.AllLangByID[langID]
	if hole == nil || lang == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if hole.Experiment != 0 || lang.Experiment != 0 {
		w.Write([]byte("[]"))
		return
	}

	otherScoring := "bytes"
	if scoring == "bytes" {
		otherScoring = "chars"
	}

	var followLimit, userID int
	if golfer := session.Golfer(r); golfer != nil {
		followLimit = golfer.FollowLimit()
		userID = golfer.ID
	}

	type entry struct {
		Bytes      *int `json:"bytes"`
		BytesChars *int `json:"bytes_chars"`
		Chars      *int `json:"chars"`
		CharsBytes *int `json:"chars_bytes"`
		Golfer     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"golfer"`
		Me   bool `json:"me"`
		Rank int  `json:"rank"`
	}

	sqlWhere, sqlLimit := "true", "$6"
	switch view {
	case "me":
		sqlWhere = "row > COALESCE((SELECT row FROM ranks WHERE me), 0) - $6"
		sqlLimit = "$6 * 2"
	case "following":
		sqlWhere = fmt.Sprint("user_id = ANY(following($1, ", followLimit, "))")
	}

	// We don't use the rankings view as we want instant updates upon solution
	// submit, therefore we skip scoring to keep it fast.
	rows, err := session.Database(r).Query(
		`WITH ranks AS (
		    SELECT ROW_NUMBER() OVER (ORDER BY `+scoring+`, submitted) row,
		           RANK()       OVER (ORDER BY `+scoring+`),
		           user_id,
		           `+scoring+`,
		           `+otherScoring+` `+scoring+`_`+otherScoring+`,
		           user_id = $1 me
		      FROM solutions
		     WHERE hole = $2
		       AND lang = $3
		       AND scoring = $4
		       AND NOT failing
		), other_scoring AS (
		    SELECT user_id,
		           `+otherScoring+`,
		           `+scoring+` `+otherScoring+`_`+scoring+`
		      FROM solutions
		     WHERE hole = $2
		       AND lang = $3
		       AND scoring = $5
		       AND NOT failing
		)   SELECT bytes, bytes_chars, chars, chars_bytes, id, login, me, rank
		      FROM ranks
		      JOIN users ON id = user_id
		 LEFT JOIN other_scoring USING(user_id)
		     WHERE `+sqlWhere+`
		  ORDER BY row
		     LIMIT `+sqlLimit,
		userID,
		holeID,
		langID,
		scoring,
		otherScoring,
		limit,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	entries := make([]entry, 0, limit)
	for rows.Next() {
		var e entry

		if err := rows.Scan(
			&e.Bytes,
			&e.BytesChars,
			&e.Chars,
			&e.CharsBytes,
			&e.Golfer.ID,
			&e.Golfer.Name,
			&e.Me,
			&e.Rank,
		); err != nil {
			panic(err)
		}

		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	// Trim the rows to limit with "me" as close to the middle as possible.
	// TODO It would simplify everything if we could fold this into the SQL.
	length := len(entries)
	if view == "me" && length > limit {
		me := slices.IndexFunc(entries, func(e entry) bool { return e.Me })
		// Before: me entries, then "me" entry, then len(entries)-me-1 entries
		// 	with me <= limit; len(entries) <= 2*limit; len(entries)-me-1 <= limit-1
		if me <= limit/2 {
			entries = entries[:limit]
		} else if me >= length-(limit+1)/2 {
			// Impossible case?
			entries = entries[length-limit:]
		} else {
			entries = entries[me-limit/2 : me+(limit+1)/2]
		}
	}

	encodeJSON(w, entries)
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

	// Only sponsors can create or update notes, the can still read & delete.
	golfer := session.Golfer(r)
	if !(golfer.Admin || golfer.Sponsor) {
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
