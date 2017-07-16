package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error

	if db, err = sql.Open("postgres", os.Getenv("DB_DSN")); err != nil {
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

func getSolutionCode(userID int, hole, lang string) (code string) {
	if err := db.QueryRow(
		`SELECT code
		   FROM solutions
		  WHERE user_id = $1 AND hole = $2 AND lang = $3`,
		userID, hole, lang,
	).Scan(&code); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return
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

func printLeaderboards(w http.ResponseWriter) {
	rows, err := db.Query(
		`SELECT hole,
		        CONCAT(
		            '<tr><td>',
		            row_number,
		            '<td>',
		            strokes,
		            '<td>',
		            CASE lang
		            WHEN 'javascript' THEN 'JavaScript'
		            WHEN 'perl6'      THEN 'Perl 6'
		            WHEN 'php'        THEN 'PHP'
		            ELSE INITCAP(CAST(lang AS TEXT))
		            END,
		            '<td><img src="//avatars.githubusercontent.com/',
		            login,
		            '?size=20"><a href="u/',
		            login,
		            '">',
		            login,
		            '</a>'
		        )
		   FROM (
		     SELECT ROW_NUMBER() OVER (
		                PARTITION BY hole ORDER BY LENGTH(code), submitted
		            ),
		            hole,
		            login,
		            lang,
		            LENGTH(code) strokes
		       FROM solutions
		       JOIN users ON user_id = id
		        ) x
		  WHERE row_number < 6`,
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
