package routes

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/buildkite/terminal-to-html/v3"
	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/discord"
	Golfer "github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

func StripBrainfuck(code string) string {
	reg, _ := regexp.Compile("[^<>\\+\\-\\[\\]\\.\\,]")
	return reg.ReplaceAllString(code, "")
}

func StripCode(code, lang, hole string) string {
	if hole == "quine" {
		return code
	}
	if lang == "brainfuck" {
		return StripBrainfuck(code)
	}
	// TODO other languages can be added here
	return code
}

// Solution serves POST /solution
func Solution(w http.ResponseWriter, r *http.Request) {
	var in struct{ Code, Hole, Lang string }

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	strippedCode := StripCode(in.Code, in.Lang, in.Hole)

	_, experimental := config.ExpHoleByID[in.Hole]
	if !experimental && config.HoleByID[in.Hole] == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, ok := config.LangByID[in.Lang]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db := session.Database(r)
	golfer := session.Golfer(r)

	// 128 KiB, >= because arguments needs a null termination.
	if len(in.Code) >= 128*1024 {
		if golfer != nil {
			golfer.Earn(db, "tl-dr")
		}

		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	score := hole.Play(r.Context(), in.Hole, in.Lang, in.Code)

	if score.Timeout && golfer != nil {
		golfer.Earn(db, "slowcoach")
	}

	out := struct {
		Argv           []string
		Cheevos        []*config.Cheevo
		Err, Exp, Out  string
		ExitCode       int
		Pass, LoggedIn bool
		RankUpdates    []Golfer.RankUpdate
		Took           time.Duration
	}{
		Argv:     score.Args,
		Cheevos:  []*config.Cheevo{},
		Err:      string(terminal.Render(score.Stderr)),
		ExitCode: score.ExitCode,
		Exp:      score.Answer,
		LoggedIn: golfer != nil,
		Out:      string(score.Stdout),
		Pass:     score.Pass,
		RankUpdates: []Golfer.RankUpdate{
			{Scoring: "bytes"},
			{Scoring: "chars"},
		},
		Took: score.Took,
	}

	if out.Pass && golfer != nil && !experimental {
		var cheevos []string
		if err := db.QueryRowContext(
			r.Context(),
			`SELECT earned,
			        old_bytes_joint, old_bytes_rank, old_bytes,
			        new_bytes_joint, new_bytes_rank, new_bytes,
			        beat_bytes,
			        old_chars_joint, old_chars_rank, old_chars,
			        new_chars_joint, new_chars_rank, new_chars,
			        beat_chars
			   FROM save_solution(
			            bytes   := CASE WHEN $3 = 'assembly'::lang
			                            THEN $5
			                            ELSE octet_length($6)
			                            END,
			            chars   := CASE WHEN $3 = 'assembly'::lang
			                            THEN NULL
			                            ELSE char_length($6)
			                            END,
			            code    := $1,
			            hole    := $2,
			            lang    := $3,
			            user_id := $4
			        )`,
			in.Code, in.Hole, in.Lang, golfer.ID, score.ASMBytes, strippedCode,
		).Scan(
			pq.Array(&cheevos),
			&out.RankUpdates[0].From.Joint,
			&out.RankUpdates[0].From.Rank,
			&out.RankUpdates[0].From.Strokes,
			&out.RankUpdates[0].To.Joint,
			&out.RankUpdates[0].To.Rank,
			&out.RankUpdates[0].To.Strokes,
			&out.RankUpdates[0].Beat,
			&out.RankUpdates[1].From.Joint,
			&out.RankUpdates[1].From.Rank,
			&out.RankUpdates[1].From.Strokes,
			&out.RankUpdates[1].To.Joint,
			&out.RankUpdates[1].To.Rank,
			&out.RankUpdates[1].To.Strokes,
			&out.RankUpdates[1].Beat,
		); err != nil {
			panic(err)
		}

		for _, cheevo := range cheevos {
			out.Cheevos = append(out.Cheevos, config.CheevoByID[cheevo])
		}

		recordUpdates := make([]Golfer.RankUpdate, 0, 2)

		for _, rank := range out.RankUpdates {
			if rank.From.Strokes.Int64 == rank.To.Strokes.Int64 {
				continue
			}

			// This keeps track of which updates (if any) represent new records
			if !rank.To.Joint.Bool && rank.To.Rank.Int64 == 1 {
				recordUpdates = append(recordUpdates, rank)
			}
		}

		// If any of the updates are record breakers, announce them on Discord
		if len(recordUpdates) > 0 {
			go discord.LogNewRecord(
				golfer, config.HoleByID[in.Hole], config.LangByID[in.Lang],
				recordUpdates, db,
			)
		}

		// TODO Use the golfer's timezone from /settings.
		var (
			now   = time.Now().UTC()
			month = now.Month()
			day   = now.Day()
		)

		if month == time.October && day == 2 {
			golfer.Earn(db, "happy-birthday-code-golf")
		}

		switch in.Hole {
		case "12-days-of-christmas":
			if (month == time.December && day >= 25) || (month == time.January && day <= 5) {
				golfer.Earn(db, "twelvetide")
			}
		case "star-wars-opening-crawl":
			if month == time.May && day == 4 {
				golfer.Earn(db, "may-the-4ᵗʰ-be-with-you")
			}
		case "united-states":
			if month == time.July && day == 4 {
				golfer.Earn(db, "independence-day")
			}
		case "vampire-numbers":
			if month == time.October && day == 31 {
				golfer.Earn(db, "vampire-byte")
			}
		case "π":
			if month == time.March && day == 14 {
				golfer.Earn(db, "pi-day")
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
