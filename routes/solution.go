package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/buildkite/terminal"
	"github.com/code-golf/code-golf/discord"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/pretty"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/guregu/null.v4"
)

// Solution serves POST /solution
func Solution(w http.ResponseWriter, r *http.Request) {
	var in struct{ Code, Hole, Lang string }

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	db := session.Database(r)
	golfer := session.Golfer(r)

	if len(in.Code) > 409_600 {
		if golfer != nil {
			awardTrophy(db, golfer.ID, "tl-dr")
		}

		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	score := hole.Play(r.Context(), in.Hole, in.Lang, in.Code)

	if score.Timeout && golfer != nil {
		awardTrophy(db, golfer.ID, "slowcoach")
	}

	type rankUpdate struct {
		Scoring  string
		From, To struct {
			Joint         null.Bool
			Rank, Strokes null.Int
		}
	}

	out := struct {
		Argv                []string
		Diff, Err, Exp, Out string
		Pass, LoggedIn      bool
		RankUpdates         []rankUpdate
		Took                time.Duration
		Trophies            []string
	}{
		Argv: score.Args,
		Err:  string(terminal.Render(score.Stderr)),
		Exp:  score.Answer,
		Out:  string(score.Stdout),
		Pass: score.Pass,
		RankUpdates: []rankUpdate{
			{Scoring: "bytes"},
			{Scoring: "chars"},
		},
		LoggedIn: golfer != nil,
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
	if out.Pass && golfer != nil && !experimental {
		if err := db.QueryRowContext(
			r.Context(),
			`SELECT earned,
			        old_bytes_joint, old_bytes_rank, old_bytes,
			        new_bytes_joint, new_bytes_rank, new_bytes,
			        old_chars_joint, old_chars_rank, old_chars,
			        new_chars_joint, new_chars_rank, new_chars
			   FROM save_solution(code := $1, hole := $2, lang := $3, user_id := $4)`,
			in.Code, in.Hole, in.Lang, golfer.ID,
		).Scan(
			pq.Array(&out.Trophies),
			&out.RankUpdates[0].From.Joint,
			&out.RankUpdates[0].From.Rank,
			&out.RankUpdates[0].From.Strokes,
			&out.RankUpdates[0].To.Joint,
			&out.RankUpdates[0].To.Rank,
			&out.RankUpdates[0].To.Strokes,
			&out.RankUpdates[1].From.Joint,
			&out.RankUpdates[1].From.Rank,
			&out.RankUpdates[1].From.Strokes,
			&out.RankUpdates[1].To.Joint,
			&out.RankUpdates[1].To.Rank,
			&out.RankUpdates[1].To.Strokes,
		); err != nil {
			panic(err)
		}

		recordScorings := ""

		for _, rank := range out.RankUpdates {
			if rank.From.Strokes.Int64 == rank.To.Strokes.Int64 {
				continue
			}

			var fromJoint, toJoint string
			if rank.From.Joint.Bool {
				fromJoint = "joint "
			}
			if rank.To.Joint.Bool {
				toJoint = "joint "
			}

			log.Printf(
				"%s: %s/%s/%s %d (%s%d%s) → %d (%s%d%s)\n",
				golfer.Name,
				in.Hole,
				in.Lang,
				rank.Scoring,
				rank.From.Strokes.Int64,
				fromJoint,
				rank.From.Rank.Int64,
				pretty.Ordinal(int(rank.From.Rank.Int64)),
				rank.To.Strokes.Int64,
				toJoint,
				rank.To.Rank.Int64,
				pretty.Ordinal(int(rank.To.Rank.Int64)),
			)

			if !rank.To.Joint.Bool && rank.To.Rank.Int64 == 1 {
				recordScorings += rank.Scoring
			}
		}

		if recordScorings != "" {
			discord.LogNewRecord(golfer, in.Hole, in.Lang, recordScorings, out.RankUpdates[0].To.Strokes.Int64, out.RankUpdates[1].To.Strokes.Int64)
		}

		// TODO Use the golfer's timezone from /settings.
		var (
			now   = time.Now().UTC()
			month = now.Month()
			day   = now.Day()
		)

		if month == time.October && day == 2 {
			awardTrophy(db, golfer.ID, "happy-birthday-code-golf")
		}

		switch in.Hole {
		case "12-days-of-christmas":
			if (month == time.December && day >= 25) || (month == time.January && day <= 5) {
				awardTrophy(db, golfer.ID, "twelvetide")
			}
		case "united-states":
			if month == time.July && day == 4 {
				awardTrophy(db, golfer.ID, "independence-day")
			}
		case "vampire-numbers":
			if month == time.October && day == 31 {
				awardTrophy(db, golfer.ID, "vampire-byte")
			}
		case "π":
			if month == time.March && day == 14 {
				awardTrophy(db, golfer.ID, "pi-day")
			}
		}

		if queryBool(
			db,
			`SELECT COUNT(DISTINCT lang) = ARRAY_LENGTH(ENUM_RANGE(NULL::lang), 1)
			   FROM solutions
			  WHERE NOT failing AND user_id = $1`,
			golfer.ID,
		) {
			awardTrophy(db, golfer.ID, "polyglot")
		}

		// FIXME Each one of these queries takes 50ms!
		if queryBool(
			db,
			"SELECT chars_points > 9000 FROM chars_points WHERE user_id = $1",
			golfer.ID,
		) || queryBool(
			db,
			"SELECT bytes_points > 9000 FROM bytes_points WHERE user_id = $1",
			golfer.ID,
		) {
			awardTrophy(db, golfer.ID, "its-over-9000")
		}

		// COUNT(*) = 4 because langs x (bytes, chars)
		switch in.Lang {
		case "java", "javascript":
			if queryBool(
				db,
				`SELECT COUNT(*) = 4
				   FROM solutions
				  WHERE NOT failing
				    AND hole = $1
				    AND lang IN ('java', 'javascript')
				    AND user_id = $2`,
				in.Hole,
				golfer.ID,
			) {
				awardTrophy(db, golfer.ID, "caffeinated")
			}
		case "perl", "raku":
			if queryBool(
				db,
				`SELECT COUNT(*) = 4
				   FROM solutions
				  WHERE NOT failing
				    AND hole = $1
				    AND lang IN ('perl', 'raku')
				    AND user_id = $2`,
				in.Hole,
				golfer.ID,
			) {
				awardTrophy(db, golfer.ID, "tim-toady")
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
