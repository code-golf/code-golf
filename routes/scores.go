package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/code-golf/code-golf/cookie"
)

func scoresMini(w http.ResponseWriter, r *http.Request) {
	userID, _ := cookie.Read(r)

	var json []byte

	if err := db(r).QueryRow(
		`WITH leaderboard AS (
		    SELECT ROW_NUMBER() OVER (ORDER BY LENGTH(code), submitted),
		           RANK()       OVER (ORDER BY LENGTH(code)),
		           user_id,
		           LENGTH(code) strokes,
		           user_id = $1 me
		      FROM solutions
		     WHERE hole = $2
		       AND lang = $3
		       AND NOT failing
		), mini_leaderboard AS (
		    SELECT rank,
		           login,
		           me,
		           strokes
		      FROM leaderboard
		      JOIN users on user_id = id
		     WHERE row_number >
		           COALESCE((SELECT row_number - 4 FROM leaderboard WHERE me), 0)
		  ORDER BY row_number
		     LIMIT 7
		) SELECT COALESCE(JSON_AGG(mini_leaderboard), '[]') FROM mini_leaderboard`,
		userID,
		param(r, "hole"),
		param(r, "lang"),
	).Scan(&json); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func scoresAll(w http.ResponseWriter, r *http.Request) {
	var json []byte

	if err := db(r).QueryRow(
		`WITH solution_lengths AS (
		    SELECT hole,
		           lang,
		           login,
		           LENGTH(code) strokes,
		           submitted
		      FROM solutions
		      JOIN users on user_id = id
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

// Scores serves GET /scores/{hole}/{lang}/{suffix}
func Scores(w http.ResponseWriter, r *http.Request) {
	holeID := param(r, "hole")
	langID := param(r, "lang")

	if _, ok := holeByID[holeID]; holeID != "all-holes" && !ok {
		render(w, r, http.StatusNotFound, "404", "", nil)
		return
	}

	if _, ok := langByID[langID]; langID != "all-langs" && !ok {
		render(w, r, http.StatusNotFound, "404", "", nil)
		return
	}

	type Score struct {
		Holes, Points, Rank, Strokes int
		Lang                         Lang
		Login                        string
		Submitted                    time.Time
	}

	data := struct {
		HoleID, LangID string
		Holes          []Hole
		Langs          []Lang
		Next, Prev     int
		Scores         []Score
	}{
		HoleID: holeID,
		Holes:  holes,
		LangID: langID,
		Langs:  langs,
	}

	page := 1

	if suffix := param(r, "suffix"); suffix != "" {
		if suffix == "mini" {
			scoresMini(w, r)
			return
		}

		if suffix == "all" {
			scoresAll(w, r)
			return
		}

		page, _ = strconv.Atoi(suffix)

		if page < 1 {
			render(w, r, http.StatusNotFound, "404", "", nil)
			return
		}

		if page == 1 {
			http.Redirect(w, r, "/scores/"+holeID+"/"+langID, http.StatusMovedPermanently)
			return
		}
	}

	if page != 1 {
		data.Prev = page - 1
	}

	var distinct, table, title string

	if holeID == "all-holes" {
		distinct = "DISTINCT ON (hole, user_id)"
		table = "summed_leaderboard"
	} else {
		table = "scored_leaderboard"
	}

	rows, err := db(r).Query(
		`WITH leaderboard AS (
		  SELECT `+distinct+`
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id,
		         lang
		    FROM solutions
		   WHERE NOT failing
		     AND $1 IN ('all-holes', hole::text)
		     AND $2 IN ('all-langs', lang::text)
		ORDER BY hole, user_id, LENGTH(code), submitted
		), scored_leaderboard AS (
		  SELECT hole,
		         1 holes,
		         ROUND(
		             (COUNT(*) OVER (PARTITION BY hole) -
		                RANK() OVER (PARTITION BY hole ORDER BY strokes) + 1)
		             * (1000.0 / COUNT(*) OVER (PARTITION BY hole))
		         ) points,
		         strokes,
		         submitted,
		         user_id,
		         lang
		    FROM leaderboard
		), summed_leaderboard AS (
		  SELECT user_id,
		         COUNT(*)       holes,
		         '' lang,
		         SUM(points)    points,
		         SUM(strokes)   strokes,
		         MAX(submitted) submitted
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT holes,
		         lang,
		         login,
		         points,
		         RANK() OVER (ORDER BY points DESC, strokes),
		         strokes,
		         submitted
		    FROM `+table+`
		    JOIN users on user_id = id
		ORDER BY points DESC, strokes, submitted
		   LIMIT 101
		  OFFSET $3`,
		holeID,
		langID,
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
			&score.Holes,
			&langID,
			&score.Login,
			&score.Points,
			&score.Rank,
			&score.Strokes,
			&score.Submitted,
		); err != nil {
			panic(err)
		}

		score.Lang = langByID[langID]

		data.Scores = append(data.Scores, score)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if holeID == "all-holes" {
		title = "All Holes"
	} else {
		title = holeByID[holeID].Name
	}

	title += " in "

	if langID == "all-langs" {
		title += "All Langs"
	} else {
		title += langByID[langID].Name
	}

	render(w, r, http.StatusOK, "scores", title, data)
}
