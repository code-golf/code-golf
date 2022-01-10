package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// ScoresMini serves GET /scores/{hole}/{lang}/{scoring}/mini
func ScoresMini(w http.ResponseWriter, r *http.Request) {
	if config.HoleByID[param(r, "hole")] == nil ||
		config.LangByID[param(r, "lang")] == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var userID int
	if golfer := session.Golfer(r); golfer != nil {
		userID = golfer.ID
	}

	scoring := param(r, "scoring")
	var otherScoring string
	switch scoring {
	case "bytes":
		otherScoring = "chars"
	case "chars":
		otherScoring = "bytes"
	default:
		NotFound(w, r)
		return
	}

	var json []byte

	limit := 7
	if r.FormValue("ng") == "1" {
		limit = 15
	}

	if err := session.Database(r).QueryRow(
		`WITH leaderboard AS (
		    SELECT ROW_NUMBER() OVER (ORDER BY `+scoring+`, submitted),
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
		), mini_leaderboard AS (
		    SELECT rank,
		           login,
		           me,
		           chars,
		           bytes,
		           chars_bytes,
		           bytes_chars
		      FROM leaderboard
		      JOIN users on user_id = id
		 LEFT JOIN other_scoring ON leaderboard.user_id = other_scoring.user_id
		     WHERE row_number >
		           COALESCE((SELECT row_number - $6 FROM leaderboard WHERE me), 0)
		  ORDER BY row_number
		     LIMIT $7
		) SELECT COALESCE(JSON_AGG(mini_leaderboard), '[]') FROM mini_leaderboard`,
		userID,
		param(r, "hole"),
		param(r, "lang"),
		scoring,
		otherScoring,
		limit/2+1,
		limit,
	).Scan(&json); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Scores serves GET /scores/{hole}/{lang}/all
func ScoresAll(w http.ResponseWriter, r *http.Request) {
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

// Scores redirects GET /scores/{hole}/{lang}/{scoring}
func Scores(w http.ResponseWriter, r *http.Request) {
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
