package main

import (
	"database/sql"
	"io"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error

	if db, err = sql.Open("postgres", ""); err != nil {
		panic(err)
	}
}

func addSolution(userID int, hole, lang, code string) {
	// Update the code if it's the same characters or less, but only update
	// the submitted time if the solution is shorter. This avoids a user
	// moving down the leaderboard by matching their personal best.
	if _, err := db.Exec(`
	    INSERT INTO solutions
	         VALUES (NOW(), $1, $2, $3, $4)
	    ON CONFLICT ON CONSTRAINT solutions_pkey
	  DO UPDATE SET submitted = CASE
	                    WHEN LENGTH($4) < LENGTH(solutions.code)
	                    THEN NOW()
	                    ELSE solutions.submitted
	                END,
	                code = CASE
	                    WHEN LENGTH($4) > LENGTH(solutions.code)
	                    THEN solutions.code
	                    ELSE $4
	                END
	`, userID, hole, lang, code); err != nil {
		panic(err)
	}
}

func getUser(login string) bool {
	var one int

	if err := db.QueryRow(
		"SELECT 1 FROM users WHERE login = $1", login,
	).Scan(&one); err != nil && err != sql.ErrNoRows {
		panic(err)
	} else {
		return err != sql.ErrNoRows
	}
}

func getUserSolutions(userID int, hole string) map[string]string {
	rows, err := db.Query(
		`SELECT code, lang
		   FROM solutions
		  WHERE user_id = $1 AND hole = $2`,
		userID, hole,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	solutions := make(map[string]string)

	for rows.Next() {
		var code, lang string

		if err := rows.Scan(&code, &lang); err != nil {
			panic(err)
		}

		solutions[lang] = code
	}

	return solutions
}

func addUser(id int, login string) {
	if _, err := db.Exec(
		`INSERT INTO users VALUES($1, $2)
		 ON CONFLICT(id) DO UPDATE SET login = $2`,
		id, login,
	); err != nil {
		panic(err)
	}
}

func printLeaderboards(w io.WriteCloser, id int) {
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
		             '?s=26"><a href="u/',
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
		         WHEN 'odious-numbers'           THEN 4
		         WHEN 'pascals-triangle'         THEN 5
		         WHEN 'pernicious-numbers'       THEN 6
		         WHEN 'prime-numbers'            THEN 7
		         WHEN '99-bottles-of-beer'       THEN 8
		         WHEN 'seven-segment'            THEN 9
		         WHEN 'sierpiński-triangle'      THEN 10
		         WHEN 'π'                        THEN 11
		         WHEN 'e'                        THEN 12
		         WHEN 'arabic-to-roman-numerals' THEN 13
		         WHEN 'spelling-numbers'         THEN 14
		         END, rank, submitted`,
		id,
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

func printScores(w io.WriteCloser, hole, lang string, userID int) {
	w.Write([]byte(
		"<script async src=" + jsScoresPath +
			"></script><article id=scores><select id=hole>",
	))

	for _, v := range [][]string{
		{"all", "All Holes"},
		{"99-bottles-of-beer", "99 Bottles of Beer"},
		{"arabic-to-roman-numerals", "Arabic to Roman"},
		{"e", "e"},
		{"emirp-numbers", "Emirp Numbers"},
		{"evil-numbers", "Evil Numbers"},
		{"fibonacci", "Fibonacci"},
		{"fizz-buzz", "Fizz Buzz"},
		{"odious-numbers", "Odious Numbers"},
		{"pascals-triangle", "Pascal's Triangle"},
		{"pernicious-numbers", "Pernicious Numbers"},
		{"prime-numbers", "Prime Numbers"},
		{"sierpiński-triangle", "Sierpiński Triangle"},
		{"seven-segment", "Seven Segment"},
		{"spelling-numbers", "Spelling Numbers"},
		{"π", "π"},
	} {
		w.Write([]byte("<option "))
		if hole == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select><select id=lang>"))

	for _, v := range [][]string{
		{"all", "All Langs"},
		{"javascript", "JavaScript"},
		{"perl", "Perl"},
		{"perl6", "Perl 6"},
		{"php", "PHP"},
		{"python", "Python"},
		{"ruby", "Ruby"},
	} {
		w.Write([]byte("<option "))
		if lang == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select><table>"))

	var concat, where string

	if hole != "" {
		where += " AND hole = '" + hole + "'"
		concat = "strokes, '<td>(', score, ' point', CASE WHEN score"
	} else {
		concat = "score, '<td>(', count, ' hole', CASE WHEN count"
	}

	if lang != "" {
		where += " AND lang = '" + lang + "'"
	}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id
		    FROM solutions
		   WHERE true`+where+`
		ORDER BY hole, user_id, LENGTH(code), submitted
		), scored_leaderboard AS (
		  SELECT hole,
		         ROUND(
		             (
		                 (SELECT COUNT(distinct user_id) FROM solutions WHERE hole = l.hole`+where+`)
		                 -
		                 RANK() OVER (PARTITION BY hole ORDER BY strokes)
		                 +
		                 1
		             )
		             *
		             (
		                 100.0
		                 /
		                 (SELECT COUNT(distinct user_id) FROM solutions WHERE hole = l.hole`+where+`)
		             )
		         ) score,
		         strokes,
		         submitted,
		         user_id
		    FROM leaderboard l
		), summed_leaderboard AS (
		  SELECT user_id,
		         SUM(strokes) strokes,
		         SUM(score)   score,
		         COUNT(*),
		         MAX(submitted)
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT CONCAT(
		             '<tr',
		             CASE WHEN user_id = $1 THEN ' class=me' END,
		             '><td>',
		             RANK() OVER (ORDER BY score DESC),
		             '<td><img src="//avatars.githubusercontent.com/',
		             login,
		             '?s=26"><a href="/u/',
		             login,
		             '">',
		             login,
		             '</a><td>',
		             `+concat+` > 1 THEN 's' END,
		             ')<td>',
		             TO_CHAR(max, 'YYYY-MM-DD<span> HH24:MI:SS</span>')
		         )
		    FROM summed_leaderboard
		    JOIN users on user_id = id
		ORDER BY score DESC, max`,
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
