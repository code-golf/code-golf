package routes

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/code-golf/code-golf/session"
)

func getOtherScoring(scoring string) string {
	if scoring == "chars" {
		return "bytes"
	}

	return "chars"
}

func scoresMini(w http.ResponseWriter, r *http.Request, scoring string) {
	var userID int
	if golfer := session.Golfer(r); golfer != nil {
		userID = golfer.ID
	}

	otherScoring := getOtherScoring(scoring)

	var json []byte

	if err := session.Database(r).QueryRow(
		`WITH leaderboard AS (
		    SELECT ROW_NUMBER() OVER (ORDER BY `+scoring+`, submitted),
		           RANK()       OVER (ORDER BY `+scoring+`),
		           user_id,
		           `+scoring+`,
		           `+otherScoring+` `+scoring+`_`+otherScoring+`,
		           user_id = $1 me
		      FROM solutions
		      JOIN code ON code_id = id
		     WHERE hole = $2
		       AND lang = $3
		       AND scoring = $4
		       AND NOT failing
		), other_scoring AS (
		    SELECT user_id,
		           `+otherScoring+`,
		           `+scoring+` `+otherScoring+`_`+scoring+`
		      FROM solutions
		      JOIN code ON code_id = id
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
		           COALESCE((SELECT row_number - 4 FROM leaderboard WHERE me), 0)
		  ORDER BY row_number
		     LIMIT 7
		) SELECT COALESCE(JSON_AGG(mini_leaderboard), '[]') FROM mini_leaderboard`,
		userID,
		param(r, "hole"),
		param(r, "lang"),
		scoring,
		otherScoring,
	).Scan(&json); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func scoresAll(w http.ResponseWriter, r *http.Request) {
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
		      JOIN code  ON code_id = code.id
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

// Scores serves GET /scores/{hole}/{lang}/{scoring}/{suffix}
func Scores(w http.ResponseWriter, r *http.Request) {
	holeID := param(r, "hole")
	langID := param(r, "lang")

	if _, ok := hole.ByID[holeID]; holeID != "all-holes" && !ok {
		NotFound(w, r)
		return
	}

	// Redirect legacy name for Raku.
	if langID == "perl6" {
		http.Redirect(
			w, r,
			strings.Replace(r.RequestURI, "perl6", "raku", 1),
			http.StatusPermanentRedirect,
		)
		return
	}

	if _, ok := lang.ByID[langID]; langID != "all-langs" && !ok {
		NotFound(w, r)
		return
	}

	scoring := param(r, "scoring")
	suffix := param(r, "suffix")

	if scoring != "chars" && scoring != "bytes" && scoring != "all" && suffix == "" {
		if scoring == "" {
			http.Redirect(w, r, "/scores/"+holeID+"/"+langID+"/chars", http.StatusPermanentRedirect)
			return
		}

		http.Redirect(w, r, "/scores/"+holeID+"/"+langID+"/chars/"+scoring, http.StatusPermanentRedirect)
		return
	}

	if scoring == "all" {
		if suffix != "" {
			NotFound(w, r)
			return
		}

		scoresAll(w, r)
		return
	}

	type Score struct {
		Bytes, Chars             int
		BytesChars, CharsBytes   int
		BytesPoints, CharsPoints int
		Country, Login           string
		Holes, Rank              int
		Lang                     lang.Lang
		Submitted                time.Time
	}

	data := struct {
		HoleID, LangID, Scoring string
		Holes                   []hole.Hole
		Langs                   []lang.Lang
		Next, Prev              int
		Scores                  []Score
	}{
		HoleID:  holeID,
		Holes:   hole.List,
		LangID:  langID,
		Langs:   lang.List,
		Scoring: scoring,
	}

	page := 1

	if suffix != "" {
		if suffix == "mini" {
			scoresMini(w, r, scoring)
			return
		}

		page, _ = strconv.Atoi(suffix)

		if page < 1 {
			NotFound(w, r)
			return
		}

		if page == 1 {
			http.Redirect(w, r, "/scores/"+holeID+"/"+langID+"/"+scoring, http.StatusPermanentRedirect)
			return
		}
	}

	if page != 1 {
		data.Prev = page - 1
	}

	var distinct, join, table, title string

	if holeID == "all-holes" {
		distinct = "DISTINCT ON (hole, user_id)"
		table = "summed_leaderboard"
	} else {
		join = "AND l.lang = o.lang"
		table = "scored_leaderboard"
	}

	otherScoring := getOtherScoring(scoring)

	rows, err := session.Database(r).Query(
		`WITH leaderboard AS (
		  SELECT `+distinct+`
		         hole,
		         submitted,
		         `+scoring+`,
		         `+otherScoring+` `+scoring+`_`+otherScoring+`,
		         user_id,
		         lang
		    FROM solutions
		    JOIN code ON code_id = id
		   WHERE NOT failing
		     AND $1 IN ('all-holes', hole::text)
		     AND $2 IN ('all-langs', lang::text)
		     AND scoring = $3
		ORDER BY hole, user_id, `+scoring+`, submitted
		), other_scoring AS (
		  SELECT `+distinct+`
		         hole,
		         `+otherScoring+`,
		         `+scoring+` `+otherScoring+`_`+scoring+`,
		         user_id,
		         lang
		    FROM solutions
		    JOIN code ON code_id = id
		   WHERE NOT failing
		     AND $1 IN ('all-holes', hole::text)
		     AND $2 IN ('all-langs', lang::text)
		     AND scoring = $4
		ORDER BY hole, user_id, `+otherScoring+`
		), scored_leaderboard AS (
		  SELECT l.hole,
		         1 holes,
		         ROUND(
		             (COUNT(*) OVER (PARTITION BY l.hole) -
		                RANK() OVER (PARTITION BY l.hole ORDER BY `+scoring+`) + 1)
		             * (1000.0 / COUNT(*) OVER (PARTITION BY l.hole))
		         ) `+scoring+`_points,
		         ROUND(
		             (COUNT(*) OVER (PARTITION BY l.hole) -
		                RANK() OVER (PARTITION BY l.hole ORDER BY `+otherScoring+`) + 1)
		             * (1000.0 / COUNT(*) OVER (PARTITION BY l.hole))
		         ) `+otherScoring+`_points,
		         COALESCE(bytes, 0) bytes,
		         COALESCE(chars, 0) chars,
		         COALESCE(bytes_chars, 0) bytes_chars,
		         COALESCE(chars_bytes, 0) chars_bytes,
		         submitted,
		         l.user_id,
		         l.lang
		    FROM leaderboard l
	   LEFT JOIN other_scoring o ON l.user_id = o.user_id AND l.hole = o.hole `+join+`
		), summed_leaderboard AS (
		  SELECT user_id,
		         COUNT(*)          holes,
		         '' lang,
		         SUM(chars_points) chars_points,
		         SUM(bytes_points) bytes_points,
		         SUM(bytes)        bytes,
		         SUM(chars)        chars,
		         SUM(bytes_chars)  bytes_chars,
		         SUM(chars_bytes)  chars_bytes,
		         MAX(submitted)    submitted
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT bytes,
		         chars,
		         bytes_chars,
		         chars_bytes,
		         COALESCE(CASE WHEN show_country THEN country END, ''),
		         holes,
		         lang,
		         login,
		         chars_points,
		         bytes_points,
		         RANK() OVER (ORDER BY `+scoring+`_points DESC, `+scoring+`),
		         submitted
		    FROM `+table+`
		    JOIN users on user_id = id
		ORDER BY `+scoring+`_points DESC, `+scoring+`, submitted
		   LIMIT 101
		  OFFSET $5`,
		holeID,
		langID,
		scoring,
		otherScoring,
		(page-1)*100,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		// We overselect by one so we can test for a next page.
		if len(data.Scores) == 100 {
			data.Next = page + 1
			continue
		}

		var langID string
		var score Score

		if err := rows.Scan(
			&score.Bytes,
			&score.Chars,
			&score.BytesChars,
			&score.CharsBytes,
			&score.Country,
			&score.Holes,
			&langID,
			&score.Login,
			&score.CharsPoints,
			&score.BytesPoints,
			&score.Rank,
			&score.Submitted,
		); err != nil {
			panic(err)
		}

		score.Lang = lang.ByID[langID]

		data.Scores = append(data.Scores, score)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if holeID == "all-holes" {
		title = "All Holes"
	} else {
		title = hole.ByID[holeID].Name
	}

	title += " in "

	if langID == "all-langs" {
		title += "All Langs"
	} else {
		title += lang.ByID[langID].Name
	}

	render(w, r, "scores", data, title)
}
