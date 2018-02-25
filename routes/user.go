package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func user(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps[0].Value

	var html []byte

	switch err := db.QueryRow(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id
		    FROM solutions
		ORDER BY hole, user_id, LENGTH(code), submitted
		), scored_leaderboard AS (
		  SELECT hole,
		         ROUND(
		             (
		                 (SELECT COUNT(distinct user_id) FROM solutions WHERE hole = l.hole)
		                 -
		                 RANK() OVER (PARTITION BY hole ORDER BY strokes)
		                 +
		                 1
		             )
		             *
		             (
		                 100.0
		                 /
		                 (SELECT COUNT(distinct user_id) FROM solutions WHERE hole = l.hole)
		             )
		         ) score,
		         strokes,
		         submitted,
		         user_id
		    FROM leaderboard l
		), summed_leaderboard AS (
		  SELECT user_id,
		         SUM(score),
		         COUNT(*),
		         MAX(submitted)
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT CONCAT(
		             login,
		             '?s=100"><h1>',
		             login,
		             '</h1><table><tr><td>',
		             TO_CHAR(sum, 'FM9,999'),
		             '<td>points<tr><td>',
		             count,
		             '<td>holes</table><hr>'
		         )
		    FROM summed_leaderboard
		    JOIN users ON id = user_id
		   WHERE login = $1`,
		user,
	).Scan(&html); err {
	case sql.ErrNoRows:
		print404(w, r)
		return
	case nil:
		printHeader(w, r, 200)
		w.Write([]byte("<link rel=stylesheet href=" + userCssPath + `><main><img src="//avatars.githubusercontent.com/`))
		w.Write(html)
	default:
		panic(err)
	}

	rows, err := db.Query(
		`WITH matrix AS (
		    SELECT user_id,
		           hole,
		           lang,
		           RANK() OVER (PARTITION BY hole, lang ORDER BY LENGTH(code))
		     FROM solutions
		) SELECT hole, lang, TO_CHAR(rank, 'FM999"<sup>"th"</sup>"')
		    FROM matrix
		    JOIN users ON id = user_id
		   WHERE login = $1
		ORDER BY CAST(hole AS text), lang`,
		user,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	w.Write([]byte("<table id=matrix><tr><th><th><th><th><th><th><th><th><th>"))

	var hole, lang, rank string

	if rows.Next() {
		if err := rows.Scan(&hole, &lang, &rank); err != nil {
			panic(err)
		}
	}

	for _, h := range holes {
		w.Write([]byte("<tr><th><a href=/" + hole + ">" + h[1] + "</a>"))

		for _, l := range []string{
			"bash", "javascript", "lua", "perl", "perl6", "php", "python", "ruby",
		} {
			w.Write([]byte("<td>"))

			if h[0] == hole && l == lang {
				w.Write([]byte("<a href=/scores/" + hole + "/" + lang + ">" + rank + "</a>"))

				if rows.Next() {
					if err := rows.Scan(&hole, &lang, &rank); err != nil {
						panic(err)
					}
				}
			} else {
				w.Write([]byte("<a href=/" + h[0] + "#" + l + "></a>"))
			}
		}
	}
}
