package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func user(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := ps[0].Value

	var html []byte
	var userID int

	switch err := db.QueryRow(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id
		    FROM solutions
		   WHERE NOT failing
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
		                 100.0
		                 /
		                 COUNT(*) OVER (PARTITION BY hole)
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
		             '<td>point',
		             CASE WHEN sum > 1 THEN 's' END,
		             '<tr><td>',
		             count,
		             '<td>hole',
		             CASE WHEN count > 1 THEN 's' END,
		             '</table><hr>'
		         ),
		         id
		    FROM summed_leaderboard
		    JOIN users ON id = user_id
		   WHERE login = $1`,
		user,
	).Scan(&html, &userID); err {
	case sql.ErrNoRows:
		print404(w, r)
		return
	case nil:
		printHeader(w, r, 200)
		w.Write([]byte("<link rel=stylesheet href=" + userCssPath + "><script defer src=" +
			timeJsPath + `></script><main><img src="//avatars.githubusercontent.com/`))
		w.Write(html)
		w.Write([]byte("<div id=trophies>"))
	default:
		panic(err)
	}

	rows, err := db.Query(
		`  SELECT CONCAT(
		              '<div',
		              CASE WHEN user_id IS NOT NULL THEN ' class=earned' END,
		              '><h2>',
		              name,
		              '</h2><p>',
		              description,
		              '<p>',
		              CASE WHEN user_id IS NULL
		                  THEN 'Not yet earned.'
		                  ELSE CONCAT(
		                      'Earned <time datetime=',
		                      TO_CHAR(earned, 'YYYY-MM-DD"T"HH24:MI:SS"Z>"FMDD Mon'),
		                      '</time>.'
		                  )
		              END,
		              '</div>'
		          )
		     FROM trophy_info i
		LEFT JOIN trophies    t ON i.trophy = t.trophy AND user_id = $1
		 ORDER BY name`,
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

	rows, err = db.Query(
		`WITH matrix AS (
		    SELECT user_id,
		           hole,
		           lang,
		           RANK() OVER (PARTITION BY hole, lang ORDER BY LENGTH(code))
		     FROM solutions
		    WHERE NOT failing
		) SELECT hole, lang, TO_CHAR(rank, 'FM999"<sup>"th"</sup>"')
		    FROM matrix
		    JOIN users ON id = user_id
		   WHERE login = $1
		ORDER BY hole, lang`,
		user,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	w.Write([]byte("</div><table id=matrix><tr><th>"))

	for _, lang := range langs {
		w.Write([]byte(`<th title="` + lang.Name + `">`))
	}

	var holeID, langID, rank string

	if rows.Next() {
		if err := rows.Scan(&holeID, &langID, &rank); err != nil {
			panic(err)
		}
	}

	for _, hole := range holes {
		w.Write([]byte("<tr><th><a href=/" + hole.ID + ">" + hole.Name + "</a>"))

		for _, lang := range langs {
			w.Write([]byte("<td>"))

			if holeID == hole.ID && langID == lang.ID {
				w.Write([]byte("<a href=/scores/" + holeID + "/" + langID + ">" + rank + "</a>"))

				if rows.Next() {
					if err := rows.Scan(&holeID, &langID, &rank); err != nil {
						panic(err)
					}
				}
			} else {
				w.Write([]byte("<a href=/" + hole.ID + "#" + lang.ID + "></a>"))
			}
		}
	}
}
