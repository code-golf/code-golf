package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header()["Strict-Transport-Security"] = []string{"max-age=31536000;includeSubDomains;preload"}

	userID := printHeader(w, r, 200)

	if userID != 0 {
		rows, err := db.Query(
			"SELECT hole, lang FROM solutions WHERE failing AND user_id = $1",
			userID,
		)

		if err != nil {
			panic(err)
		}

		defer rows.Close()

		looped := false

		for rows.Next() {
			if !looped {
				w.Write([]byte("<div id=failing>The following of your solutions have been marked as failing and no longer contribute to scoring, please update them to pass:<ul>"))
				looped = true
			}

			var holeID, langID string

			if err := rows.Scan(&holeID, &langID); err != nil {
				panic(err)
			}

			w.Write([]byte(
				"<li><a href=" + holeID + "#" + langID + ">" +
					holeByID[holeID].Name + " (" + langByID[langID].Name +
					")</a>",
			))
		}

		if looped {
			w.Write([]byte("</ul></div>"))
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}
	}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         lang,
		         LENGTH(code) strokes,
		         submitted,
		         user_id
		    FROM solutions
		   WHERE NOT failing
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
		             '<td><img src="//avatars.githubusercontent.com/',
		             login,
		             '?s=26"><a href="users/',
		             login,
		             '">',
		             login,
		             '</a><td class=',
		             lang,
		             '>',
		             strokes
		         )
		    FROM ranked_leaderboard
		    JOIN users on user_id = id
		   WHERE row_number < 6
		      OR user_id = $1
		ORDER BY CASE hole
		         WHEN 'divisors'             THEN 0
		         WHEN 'emirp-numbers'        THEN 1
		         WHEN 'evil-numbers'         THEN 2
		         WHEN 'fibonacci'            THEN 3
		         WHEN 'fizz-buzz'            THEN 4
		         WHEN 'happy-numbers'        THEN 5
		         WHEN 'niven-numbers'        THEN 6
		         WHEN 'odious-numbers'       THEN 7
		         WHEN 'pascals-triangle'     THEN 8
		         WHEN 'pernicious-numbers'   THEN 9
		         WHEN 'prime-numbers'        THEN 10
		         WHEN 'quine'                THEN 11
		         WHEN '12-days-of-christmas' THEN 12
		         WHEN '99-bottles-of-beer'   THEN 13
		         WHEN 'christmas-trees'      THEN 14
		         WHEN 'morse-decoder'        THEN 15
		         WHEN 'morse-encoder'        THEN 16
		         WHEN 'pangram-grep'         THEN 17
		         WHEN 'seven-segment'        THEN 18
		         WHEN 'sierpiński-triangle'  THEN 19
		         WHEN 'π'                    THEN 20
		         WHEN 'φ'                    THEN 21
		         WHEN '𝑒'                    THEN 22
		         WHEN 'τ'                    THEN 23
		         WHEN 'arabic-to-roman'      THEN 24
		         WHEN 'brainfuck'            THEN 25
		         WHEN 'roman-to-arabic'      THEN 26
		         WHEN 'spelling-numbers'     THEN 27
		         END, row_number`,
		userID,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	w.Write([]byte("<ul><li>Fast<li>Medium<li>Slow</ul><main id=home>"))

	var i uint8
	var prevHoleID string

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
		w.Write([]byte(prevHoleID))
		w.Write([]byte(">FULL LEADERBOARD</a></div>"))
	}

	for rows.Next() {
		var holeID string
		var row []byte

		if err := rows.Scan(&holeID, &row); err != nil {
			panic(err)
		}

		if holeID != prevHoleID {
			if prevHoleID != "" {
				printHole()
				i = 0
			}

			hole := holeByID[holeID]

			w.Write([]byte(
				"<div class=" + hole.Difficulty + "><a href=" + holeID + ">" +
					hole.Name + "</a><table class=scores>",
			))

			prevHoleID = holeID
		}

		holeRows[i] = row
		i++
	}

	printHole()

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
