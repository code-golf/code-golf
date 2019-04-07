package routes

import (
	"database/sql"
	"net/http"
	"strings"

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
	var holeID, langID string

	showDuplicates := strings.HasSuffix(r.URL.Path, "/all")

	if len(ps) == 1 {
		param := ps[0].Value

		if _, ok := holeByID[param]; ok {
			holeID = param
		} else if _, ok = langByID[param]; ok {
			langID = param
		} else {
			print404(w, r)
			return
		}
	} else if len(ps) == 2 {
		holeID = ps[0].Value
		langID = ps[1].Value

		if _, ok := holeByID[holeID]; !ok {
			print404(w, r)
			return
		}

		if langID == "all" {
			langID = ""
		} else if _, ok := langByID[langID]; !ok {
			print404(w, r)
			return
		}
	}

	userID := printHeader(w, r, 200)

	w.Write([]byte(
		"<script defer src=" + scoresJsPath + "></script><script defer src=" +
			timeJsPath + "></script><main id=scores><select id=hole><option value>All Holes",
	))

	for _, hole := range holes {
		w.Write([]byte("<option "))
		if holeID == hole.ID {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + hole.ID + ">" + hole.Name))
	}

	w.Write([]byte("</select><select id=lang><option value>All Langs"))

	for _, lang := range langs {
		w.Write([]byte("<option "))
		if langID == lang.ID {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + lang.ID + ">" + lang.Name))
	}

	w.Write([]byte("</select>"))

	if holeID != "" && langID == "" {
		w.Write([]byte("<label><input type=checkbox"))
		if showDuplicates {
			w.Write([]byte(" checked"))
		}
		w.Write([]byte(">Allow multiple entries per player</label>"))
	}

	w.Write([]byte("<table class=scores>"))

	var concat, distinct, table, where string

	if holeID != "" {
		where += " AND hole = '" + holeID + "'"
		concat = "' class=', lang, '>', TO_CHAR(strokes, 'FM999,999'), '<td>(', TO_CHAR(score, 'FM99,999'), ' point', CASE WHEN score"
	} else {
		concat = "'>', TO_CHAR(score, 'FM99,999'), '<td>(', count, ' hole', CASE WHEN count"
	}

	if showDuplicates {
		table = "scored_leaderboard"
	} else {
		distinct = "DISTINCT ON (hole, user_id)"
		table = "summed_leaderboard"
	}

	if langID != "" {
		where += " AND lang IN('" + langID + "')"
	}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT `+distinct+`
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id,
		         lang
		    FROM solutions
		   WHERE NOT failing`+where+`
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
		             CASE WHEN user_id = $1 THEN ' class=me' END,
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
		userID,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

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
