package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/buildkite/terminal"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
	"github.com/pmezard/go-difflib/difflib"
)

// Solution serves POST /solution
func Solution(w http.ResponseWriter, r *http.Request) {
	var in struct{ Code, Hole, Lang string }

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var userID int
	if golfer := session.Golfer(r); golfer != nil {
		userID = golfer.ID
	}

	db := session.Database(r)

	score := hole.Play(r.Context(), in.Hole, in.Lang, in.Code)

	if score.Timeout && userID != 0 {
		awardTrophy(db, userID, "slowcoach")
	}

	out := struct {
		Argv                []string
		Diff, Err, Exp, Out string
		Pass, LoggedIn      bool
		Took                time.Duration
	}{
		Argv:     score.Args,
		Err:      string(terminal.Render(score.Stderr)),
		Exp:      score.Answer,
		Out:      string(score.Stdout),
		Pass:     score.Pass,
		LoggedIn: userID != 0,
		Took:     score.Took,
	}

	out.Diff, _ = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(out.Exp),
		B:        difflib.SplitLines(out.Out),
		Context:  3,
		FromFile: "Exp",
		ToFile:   "Out",
	})

	_, experimental := hole.ExperimentalByID[in.Hole]
	if out.Pass && userID != 0 && !experimental {
		if _, err := db.Exec(
			"INSERT INTO code (code) VALUES ($1) ON CONFLICT DO NOTHING",
			in.Code,
		); err != nil {
			panic(err)
		}

		// Update the code if it's the same length or less, but only update
		// the submitted time if the solution is shorter. This avoids a user
		// moving down the leaderboard by matching their personal best.
		scorings := []string{"chars"}
		if session.Beta(r) {
			scorings = append(scorings, "bytes")
		}

		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			panic(err)
		}

		defer tx.Rollback()

		for _, scoring := range scorings {
			if _, err := tx.Exec(
				`WITH new_code AS (
				    SELECT id, `+scoring+` FROM code WHERE code = $1
				) INSERT INTO solutions (code_id, user_id, hole, lang, scoring)
				       SELECT id, $2, $3, $4, $5 FROM new_code
				  ON CONFLICT ON CONSTRAINT solutions_pkey
				DO UPDATE SET failing = false,
				            submitted = CASE
				                WHEN solutions.failing
				                  OR (SELECT `+scoring+` FROM new_code)
				                   < (SELECT `+scoring+` FROM code WHERE id = solutions.code_id)
				                THEN excluded.submitted
				                ELSE solutions.submitted
				            END,
				              code_id = CASE
				                WHEN solutions.failing
				                  OR (SELECT `+scoring+` FROM new_code)
				                  <= (SELECT `+scoring+` FROM code WHERE id = solutions.code_id)
				                THEN excluded.code_id
				                ELSE solutions.code_id
				            END`,
				in.Code, userID, in.Hole, in.Lang, scoring,
			); err != nil {
				panic(err)
			}
		}

		if err := tx.Commit(); err != nil {
			panic(err)
		}

		awardTrophy(db, userID, "hello-world")

		// TODO Use the golfer's timezone from /settings.
		var (
			now   = time.Now().UTC()
			month = now.Month()
			day   = now.Day()
		)

		if month == time.October && day == 2 {
			awardTrophy(db, userID, "happy-birthday-code-golf")
		}

		switch in.Hole {
		case "12-days-of-christmas":
			if (month == time.December && day >= 25) || (month == time.January && day <= 5) {
				awardTrophy(db, userID, "twelvetide")
			}
		case "fizz-buzz":
			awardTrophy(db, userID, "interview-ready")
		case "united-states":
			if month == time.July && day == 4 {
				awardTrophy(db, userID, "independence-day")
			}
		case "π":
			if month == time.March && day == 14 {
				awardTrophy(db, userID, "pi-day")
			}
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
		"INSERT INTO trophies VALUES (DEFAULT, $1, $2) ON CONFLICT DO NOTHING",
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
