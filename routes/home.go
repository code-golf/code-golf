package routes

import (
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
		         RANK()       OVER (PARTITION BY hole ORDER BY strokes),
		         ROW_NUMBER() OVER (PARTITION BY hole ORDER BY strokes, submitted),
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
		             TO_CHAR(rank, 'FM999"<sup>"th"</sup>"'),
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
		   WHERE row_number < 6
		      OR user_id = $1
		ORDER BY CASE hole
		         WHEN 'emirp-numbers'       THEN 0
		         WHEN 'evil-numbers'        THEN 1
		         WHEN 'fibonacci'           THEN 2
		         WHEN 'fizz-buzz'           THEN 3
		         WHEN 'happy-numbers'       THEN 4
		         WHEN 'odious-numbers'      THEN 5
		         WHEN 'pascals-triangle'    THEN 6
		         WHEN 'pernicious-numbers'  THEN 7
		         WHEN 'prime-numbers'       THEN 8
		         WHEN 'quine'               THEN 9
		         WHEN '99-bottles-of-beer'  THEN 10
		         WHEN 'seven-segment'       THEN 11
		         WHEN 'sierpiński-triangle' THEN 12
		         WHEN 'π'                   THEN 13
		         WHEN 'φ'                   THEN 14
		         WHEN '𝑒'                   THEN 15
		         WHEN 'arabic-to-roman'     THEN 16
		         WHEN 'roman-to-arabic'     THEN 17
		         WHEN 'spelling-numbers'    THEN 18
		         END, row_number`,
		printHeader(w, r, 200),
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	w.Write([]byte("<ul><li>Fast<li>Medium<li>Slow</ul><article id=home>"))

	var i uint8
	var prevHole string

	holeRows := make([][]byte, 6)

	printHole := func() {
		if i == 6 {
			holeRows[3] = []byte("<tr><td colspan=3>…")
			holeRows[4] = holeRows[5]
			i = 5
		}

		for j := uint8(0); j < i; j++ {
			w.Write(holeRows[j])
		}

		// Fill in blank rows if we have too few rows.
		for j := i; j < 5; j++ {
			w.Write([]byte("<tr><td><td><td>"))
		}

		w.Write([]byte("</table><a href=scores/"))
		w.Write([]byte(prevHole))
		w.Write([]byte("/all>FULL LEADERBOARD</a></div>"))
	}

	for rows.Next() {
		var hole string
		var row []byte

		if err := rows.Scan(&hole, &row); err != nil {
			panic(err)
		}

		if hole != prevHole {
			if prevHole != "" {
				printHole()
				i = 0
			}

			w.Write([]byte(`<div class=`))

			switch hole {
			case "99-bottles-of-beer":
				w.Write([]byte(`Medium><a href=99-bottles-of-beer>99 Bottles of Beer`))
			case "arabic-to-roman":
				w.Write([]byte(`Slow><a href=arabic-to-roman>Arabic to Roman`))
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
			case "quine":
				w.Write([]byte(`Fast><a href=quine>Quine`))
			case "roman-to-arabic":
				w.Write([]byte(`Slow><a href=roman-to-arabic>Roman to Arabic`))
			case "seven-segment":
				w.Write([]byte(`Medium><a href=seven-segment>Seven Segment`))
			case "sierpiński-triangle":
				w.Write([]byte(`Medium><a href=sierpiński-triangle>Sierpiński Triangle`))
			case "spelling-numbers":
				w.Write([]byte(`Slow><a href=spelling-numbers>Spelling Numbers`))
			case "π":
				w.Write([]byte(`Medium><a href=π>π`))
			case "φ":
				w.Write([]byte(`Medium><a href=φ>φ`))
			case "𝑒":
				w.Write([]byte(`Medium><a href=𝑒>𝑒`))
			}

			w.Write([]byte(`</a><table>`))

			prevHole = hole
		}

		holeRows[i] = row
		i++
	}

	printHole()

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
