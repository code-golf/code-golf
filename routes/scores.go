package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/session"
)

// GET /scores/{hole}/{lang}/all
func scoresAllGET(w http.ResponseWriter, r *http.Request) {
	var json []byte

	if err := session.Database(r).QueryRow(
		`WITH solution_lengths AS (
		    SELECT hole,
		           lang,
		           scoring,
		           login,
		           chars,
		           bytes,
		           submitted
		      FROM solutions
		      JOIN users ON user_id = users.id
		     WHERE NOT failing
		       AND $1 IN ('all-holes', hole::text)
		       AND $2 IN ('all-langs', lang::text)
		) SELECT COALESCE(JSON_AGG(solution_lengths), '[]') FROM solution_lengths`,
		param(r, "hole"),
		param(r, "lang"),
	).Scan(&json); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// GET /scores/{hole}/{lang}/{scoring}
func scoresGET(w http.ResponseWriter, r *http.Request) {
	holeID := param(r, "hole")
	langID := param(r, "lang")
	scoring := param(r, "scoring")

	if holeID == "all-holes" {
		holeID = "all"
	}

	switch langID {
	case "all-langs":
		langID = "all"
	case "perl6":
		langID = "raku"
	}

	if scoring == "" {
		scoring = "bytes"
	}

	http.Redirect(w, r, "/rankings/holes/"+holeID+"/"+langID+"/"+scoring, http.StatusPermanentRedirect)
}
