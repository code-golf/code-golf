package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header()["Strict-Transport-Security"] = []string{"max-age=31536000;includeSubDomains;preload"}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         lang,
		         LENGTH(code) strokes,
		         submitted,
		         user_id
		    FROM solutions
		ORDER BY hole, user_id, LENGTH(code), submitted
		), ranked_leaderboard AS (
		  SELECT hole,
		         lang,
		         RANK() OVER (PARTITION BY hole ORDER BY strokes),
		         strokes,
		         submitted,
		         user_id
		    FROM leaderboard
		) SELECT hole,
		         CONCAT(
		             '<tr title="Submitted ',
		             TO_CHAR(submitted, 'YYYY-MM-DD HH24:MI:SS'),
		             '"',
		             CASE WHEN user_id = $1 THEN ' class=me' END,
		             '><td>',
		             rank,
		             '<td class=',
		             lang,
		             '>',
		             strokes,
		             '<td><img src="//avatars.githubusercontent.com/',
		             login,
		             '?s=26"><a href="users/',
		             login,
		             '">',
		             login,
		             '</a>'
		         )
		    FROM ranked_leaderboard
		    JOIN users on user_id = id
		   WHERE rank < 6
		ORDER BY CASE hole
		         WHEN 'emirp-numbers'            THEN 0
		         WHEN 'evil-numbers'             THEN 1
		         WHEN 'fibonacci'                THEN 2
		         WHEN 'fizz-buzz'                THEN 3
		         WHEN 'happy-numbers'            THEN 4
		         WHEN 'odious-numbers'           THEN 5
		         WHEN 'pascals-triangle'         THEN 6
		         WHEN 'pernicious-numbers'       THEN 7
		         WHEN 'prime-numbers'            THEN 8
		         WHEN '99-bottles-of-beer'       THEN 9
		         WHEN 'seven-segment'            THEN 10
		         WHEN 'sierpiński-triangle'      THEN 11
		         WHEN 'π'                        THEN 12
		         WHEN 'e'                        THEN 13
		         WHEN 'arabic-to-roman-numerals' THEN 14
		         WHEN 'spelling-numbers'         THEN 15
		         END, rank, submitted`,
		printHeader(w, r, 200),
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	w.Write([]byte("<ul><li>Fast<li>Medium<li>Slow</ul><article id=home>"))

	var i uint8
	var prevHole string

	for rows.Next() {
		var hole string
		var row sql.RawBytes

		if err := rows.Scan(&hole, &row); err != nil {
			panic(err)
		}

		if hole != prevHole {
			if prevHole != "" {
				for j := i; j < 5; j++ {
					w.Write([]byte("<tr><td><td><td>"))
				}

				i = 0

				w.Write([]byte("</table><a href=scores/"))
				w.Write([]byte(prevHole))
				w.Write([]byte("/all>FULL LEADERBOARD</a></div>"))
			}

			w.Write([]byte(`<div class=`))

			switch hole {
			case "99-bottles-of-beer":
				w.Write([]byte(`Medium><a href=99-bottles-of-beer>99 Bottles of Beer`))
			case "arabic-to-roman-numerals":
				w.Write([]byte(`Slow><a href=arabic-to-roman-numerals>Arabic to Roman`))
			case "e":
				w.Write([]byte(`Medium><a href=e>e`))
			case "emirp-numbers":
				w.Write([]byte(`Fast><a href=emirp-numbers>Emirp Numbers`))
			case "evil-numbers":
				w.Write([]byte(`Fast><a href=evil-numbers>Evil Numbers`))
			case "fibonacci":
				w.Write([]byte(`Fast><a href=fibonacci>Fibonacci`))
			case "fizz-buzz":
				w.Write([]byte(`Fast><a href=fizz-buzz>Fizz Buzz`))
			case "happy-numbers":
				w.Write([]byte(`Fast><a href=happy-numbers>Happy Numbers`))
			case "odious-numbers":
				w.Write([]byte(`Fast><a href=odious-numbers>Odious Numbers`))
			case "pascals-triangle":
				w.Write([]byte(`Fast><a href=pascals-triangle>Pascal's Triangle`))
			case "pernicious-numbers":
				w.Write([]byte(`Fast><a href=pernicious-numbers>Pernicious Numbers`))
			case "prime-numbers":
				w.Write([]byte(`Fast><a href=prime-numbers>Prime Numbers`))
			case "seven-segment":
				w.Write([]byte(`Medium><a href=seven-segment>Seven Segment`))
			case "sierpiński-triangle":
				w.Write([]byte(`Medium><a href=sierpiński-triangle>Sierpiński Triangle`))
			case "spelling-numbers":
				w.Write([]byte(`Slow><a href=spelling-numbers>Spelling Numbers`))
			case "π":
				w.Write([]byte(`Medium><a href=π>π`))
			}

			w.Write([]byte(`</a><table>`))

			prevHole = hole
		}

		if i < 5 {
			w.Write(row)
		}

		i++
	}

	w.Write([]byte("</table><a href=scores/"))
	w.Write([]byte(prevHole))
	w.Write([]byte("/all>FULL LEADERBOARD</a></div>"))

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
