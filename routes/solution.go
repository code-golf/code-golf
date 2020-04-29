package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/buildkite/terminal"
	"github.com/code-golf/code-golf/cookie"
	"github.com/code-golf/code-golf/hole"
	"github.com/pmezard/go-difflib/difflib"
)

// Solution serves POST /solution
func Solution(w http.ResponseWriter, r *http.Request) {
	var in struct{ Code, Hole, Lang string }

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	println(in.Code)

	userID, _ := cookie.Read(r)

	db := db(r)

	score := hole.Play(r.Context(), in.Hole, in.Lang, in.Code)

	if score.Timeout && userID != 0 {
		awardTrophy(db, userID, "slowcoach")
	}

	out := struct {
		Argv                []string
		Diff, Err, Exp, Out string
		Pass                bool
		Took                time.Duration
	}{
		Argv: score.Args,
		Err:  string(terminal.Render(score.Stderr)),
		Exp:  score.Answer,
		Out:  string(score.Stdout),
		Pass: score.Pass,
		Took: score.Took,
	}

	out.Diff, _ = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(out.Exp),
		B:        difflib.SplitLines(out.Out),
		Context:  3,
		FromFile: "Exp",
		ToFile:   "Out",
	})

	if out.Pass && userID != 0 {
		// Update the code if it's the same length or less, but only update
		// the submitted time if the solution is shorter. This avoids a user
		// moving down the leaderboard by matching their personal best.
		if _, err := db.Exec(`
		    INSERT INTO solutions
		         VALUES (NOW() AT TIME ZONE 'UTC', $1, $2, $3, $4)
		    ON CONFLICT ON CONSTRAINT solutions_pkey
		  DO UPDATE SET failing = false,
		                submitted = CASE
		                    WHEN solutions.failing OR LENGTH($4) < LENGTH(solutions.code)
		                    THEN NOW() AT TIME ZONE 'UTC'
		                    ELSE solutions.submitted
		                END,
		                code = CASE
		                    WHEN LENGTH($4) > LENGTH(solutions.code) AND NOT solutions.failing
		                    THEN solutions.code
		                    ELSE $4
		                END
		`, userID, in.Hole, in.Lang, in.Code); err != nil {
			panic(err)
		}

		awardTrophy(db, userID, "hello-world")

		if now := time.Now().UTC(); now.Day() == 2 && now.Month() == time.October {
			awardTrophy(db, userID, "happy-birthday-code-golf")
		}

		if in.Hole == "fizz-buzz" {
			awardTrophy(db, userID, "interview-ready")
		}

		if queryBool(
			db,
			`SELECT COUNT(DISTINCT hole) > 18
			   FROM solutions
			  WHERE NOT failing AND user_id = $1`,
			userID,
		) {
			awardTrophy(db, userID, "the-watering-hole")
		}

		if queryBool(
			db,
			`SELECT COUNT(DISTINCT lang) = ARRAY_LENGTH(ENUM_RANGE(NULL::lang), 1)
			   FROM solutions
			  WHERE NOT failing AND user_id = $1`,
			userID,
		) {
			awardTrophy(db, userID, "polyglot")
		}

		if queryBool(
			db,
			"SELECT points > 9000 FROM points WHERE user_id = $1",
			userID,
		) {
			awardTrophy(db, userID, "its-over-9000")
		}

		switch in.Lang {
		case "brainfuck":
			if in.Hole == "brainfuck" {
				awardTrophy(db, userID, "inception")
			}
		case "java", "javascript":
			if queryBool(
				db,
				`SELECT COUNT(*) = 2
				   FROM solutions
				  WHERE NOT failing
				    AND hole = $1
				    AND lang IN ('java', 'javascript')
				    AND user_id = $2`,
				in.Hole,
				userID,
			) {
				awardTrophy(db, userID, "caffeinated")
			}
		case "php":
			awardTrophy(db, userID, "elephpant-in-the-room")
		case "perl", "raku":
			if queryBool(
				db,
				`SELECT COUNT(*) = 2
				   FROM solutions
				  WHERE NOT failing
				    AND hole = $1
				    AND lang IN ('perl', 'raku')
				    AND user_id = $2`,
				in.Hole,
				userID,
			) {
				awardTrophy(db, userID, "tim-toady")
			}
		case "python":
			if in.Hole == "quine" {
				awardTrophy(db, userID, "ouroboros")
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(&out); err != nil {
		panic(err)
	}
}

func awardTrophy(db *sql.DB, userID int, trophy string) {
	if _, err := db.Exec(
		"INSERT INTO trophies VALUES(NOW() AT TIME ZONE 'UTC', $1, $2) ON CONFLICT DO NOTHING",
		userID,
		trophy,
	); err != nil {
		panic(err)
	}
}

func queryBool(db *sql.DB, query string, args ...interface{}) (b bool) {
	if err := db.QueryRow(query, args...).Scan(&b); err != nil {
		panic(err)
	}

	return
}
