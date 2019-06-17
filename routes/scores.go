package routes

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/JRaspass/code-golf/cookie"
	"github.com/julienschmidt/httprouter"
)

func scoresMini(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, _ := cookie.Read(r)

	rows, err := db.Query(
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
		           strokes,
		           me
		      FROM leaderboard
		      JOIN users on user_id = id
		     WHERE row_number >
		           (SELECT row_number FROM leaderboard WHERE me) - 4
		  ORDER BY row_number
		     LIMIT 7
		) SELECT CONCAT(
		    '<tr',
		    CASE WHEN me THEN ' class=me' END,
		    '><td>',
		    TO_CHAR(rank, 'FM9,999"<sup>"th"</sup>"'),
		    '<td><img src="//avatars.githubusercontent.com/',
		    login,
		    '?s=26"><td>',
		    strokes
		) FROM mini_leaderboard`,
		userID,
		ps[0].Value,
		ps[1].Value,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	w.Header().Set("Content-Type", "text/html;charset=utf8")
	w.Write([]byte("<table>"))

	for rows.Next() {
		var row sql.RawBytes

		if err := rows.Scan(&row); err != nil {
			panic(err)
		}

		w.Write(row)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func scores(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	holeID := ps[0].Value
	langID := ps[1].Value

	if _, ok := holeByID[holeID]; holeID != "all-holes" && !ok {
		Render(w, r, http.StatusNotFound, "404", "", nil)
		return
	}

	if _, ok := langByID[langID]; langID != "all-langs" && !ok {
		Render(w, r, http.StatusNotFound, "404", "", nil)
		return
	}

	data := struct {
		HoleID, LangID, ScoresJsPath, TimeJsPath string
		Holes                                    []Hole
		Langs                                    []Lang
		Table                                    template.HTML
	}{
		HoleID:       holeID,
		Holes:        holes,
		LangID:       langID,
		Langs:        langs,
		ScoresJsPath: scoresJsPath,
		TimeJsPath:   timeJsPath,
	}

	var concat, distinct, table string

	if holeID == "all-holes" {
		distinct = "DISTINCT ON (hole, user_id)"
		table = "summed_leaderboard"
		concat = "'>', TO_CHAR(score, 'FM99,999'), '<td>(', count, ' hole', CASE WHEN count"
	} else {
		table = "scored_leaderboard"
		concat = "' class=', lang, '>', TO_CHAR(strokes, 'FM999,999'), '<td>(', TO_CHAR(score, 'FM99,999'), ' point', CASE WHEN score"
	}

	userID, _ := cookie.Read(r)

	rows, err := db.Query(
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
		         ROUND(
		             (
		                 COUNT(*) OVER (PARTITION BY hole)
		                 -
		                 RANK() OVER (PARTITION BY hole ORDER BY strokes)
		                 +
		                 1
		             )
		             *
		             (
		                 1000.0
		                 /
		                 COUNT(*) OVER (PARTITION BY hole)
		             )
		         ) score,
		         strokes,
		         submitted,
		         user_id,
		         lang
		    FROM leaderboard
		), summed_leaderboard AS (
		  SELECT user_id,
		         SUM(strokes)   strokes,
		         SUM(score)     score,
		         COUNT(*),
		         MAX(submitted) submitted,
		         STRING_AGG(CONCAT(lang), '') lang
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT CONCAT(
		             '<tr',
		             CASE WHEN user_id = $3 THEN ' class=me' END,
		             '><td>',
		             TO_CHAR(
		                 RANK() OVER (ORDER BY score DESC, strokes),
		                 'FM9,999"<sup>"th"</sup>"'
		             ),
		             '<td><img src="//avatars.githubusercontent.com/',
		             login,
		             '?s=26"><a href="/users/',
		             login,
		             '">',
		             login,
		             '</a><td',
		             `+concat+` > 1 THEN 's' END,
		             ')<td><time datetime=',
		             TO_CHAR(submitted, 'YYYY-MM-DD"T"HH24:MI:SS"Z>"FMDD Mon'),
		             '</time>'
		         )
		    FROM `+table+`
		    JOIN users on user_id = id
		ORDER BY score DESC, strokes, submitted`,
		holeID,
		langID,
		userID,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var row template.HTML

		if err := rows.Scan(&row); err != nil {
			panic(err)
		}

		data.Table += row
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	Render(w, r, http.StatusOK, "scores", "Scores", data)
}
