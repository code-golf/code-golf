package main

import (
	"database/sql"
	"io"
	"math/big"
	"sort"
	"strconv"

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
		`SELECT hole,
		        CONCAT(
		            '<tr title="Submitted ',
		            TO_CHAR(submitted, 'YYYY-MM-DD HH24:MI:SS'),
		            '"',
		            CASE WHEN user_id = $1 THEN ' class=me' END,
		            '><td>',
		            place,
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
		   FROM leaderboard
		   JOIN users ON user_id = id
		  WHERE place < 6
		  ORDER BY hole, place, submitted`,
		id,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	prevHole := ""
	w.Write([]byte("<article id=home>"))

	for rows.Next() {
		var hole, row string

		if err := rows.Scan(&hole, &row); err != nil {
			panic(err)
		}

		if hole != prevHole {
			if prevHole != "" {
				w.Write([]byte("</table></div>"))
			}
			w.Write([]byte(intros[hole]))
			prevHole = hole
		}

		w.Write([]byte(row))
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func printOverallLeaderboards(w io.WriteCloser, login string) {
	rows, err := db.Query(
		`SELECT login,
		        place,
		        (SELECT COUNT(*) FROM leaderboard WHERE hole = l.hole)
		   FROM leaderboard l
		   JOIN users ON user_id = id`,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	type total struct {
		login  string
		score  *big.Rat
		iScore int
		holes  int
	}

	var totals []total

	for rows.Next() {
		var login string
		var place, count int64

		if err := rows.Scan(&login, &place, &count); err != nil {
			panic(err)
		}

		score := big.NewRat(100, count)
		score.Mul(score, new(big.Rat).SetInt64(count-place+1))

		l := len(totals)

		if l == 0 || totals[l-1].login != login {
			totals = append(totals, total{login, score, 0, 1})
		} else {
			totals[l-1].holes++
			totals[l-1].score = score.Add(score, totals[l-1].score)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	// Convert the rational numbers to integers.
	for i, v := range totals {
		score, _ := v.score.Float64()
		totals[i].iScore = int(score)
	}

	sort.Slice(totals, func(i, j int) bool { return totals[i].iScore > totals[j].iScore })

	w.Write([]byte("<article><table id=leaderboard>"))

	for i, v := range totals {
		w.Write([]byte("<tr"))

		if login == v.login {
			w.Write([]byte(" class=me"))
		}

		w.Write([]byte(
			"><td>" + strconv.Itoa(i+1) +
				`<td><img src="//avatars.githubusercontent.com/` + v.login +
				`?s=26"><a href="u/` + v.login + `">` + v.login + "</a>" +
				"<td>" + strconv.Itoa(v.iScore) +
				" <i>(" + strconv.Itoa(v.holes) + " holes)</i>",
		))
	}
}
