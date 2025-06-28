package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/discord"
	Golfer "github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /scores/{hole}/{lang}/all
func scoresAllGET(w http.ResponseWriter, r *http.Request) {
	var json []byte

	if err := session.Database(r).QueryRow(
		`WITH solution_lengths AS (
		    SELECT hole,
		           lang,
		           scoring,
		           login,
		           chars,
		           bytes,
		           submitted
		      FROM solutions
		      JOIN users ON user_id = users.id
		     WHERE NOT failing
		       AND $1 IN ('all-holes', hole::text)
		       AND $2 IN ('all-langs', lang::text)
		  ORDER BY submitted DESC
		) SELECT COALESCE(JSON_AGG(solution_lengths), '[]') FROM solution_lengths`,
		param(r, "hole"),
		param(r, "lang"),
	).Scan(&json); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// GET /scores/{hole}/{lang}/{scoring}
func scoresGET(w http.ResponseWriter, r *http.Request) {
	holeID := param(r, "hole")
	langID := param(r, "lang")
	scoring := param(r, "scoring")

	if holeID == "all-holes" {
		holeID = "all"
	}
	if langID == "all-holes" {
		langID = "all"
	}
	if scoring == "" {
		scoring = "bytes"
	}

	http.Redirect(w, r, "/rankings/holes/"+holeID+"/"+langID+"/"+scoring, http.StatusPermanentRedirect)
}

// POST /solution
func solutionPOST(w http.ResponseWriter, r *http.Request) {
	var in struct{ Code, Hole, Lang string }

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	holeObj := config.AllHoleByID[in.Hole]
	langObj := config.AllLangByID[in.Lang]
	if holeObj == nil || langObj == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	experimentalHole := holeObj.Experiment != 0
	experimentalLang := langObj.Experiment != 0
	experimental := experimentalHole || experimentalLang

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

	runs, err := hole.Play(r.Context(), holeObj, langObj, in.Code)
	if err != nil {
		panic(err)
	}

	out := struct {
		Cheevos     []config.Cheevo     `json:"cheevos"`
		LoggedIn    bool                `json:"logged_in"`
		RankUpdates []Golfer.RankUpdate `json:"rank_updates"`
		Runs        []hole.Run          `json:"runs"`
	}{
		LoggedIn: golfer != nil,
		Runs:     runs,
	}

	if golfer != nil && slices.ContainsFunc(
		out.Runs, func(r hole.Run) bool { return r.Timeout },
	) {
		if c := golfer.Earn(db, "slowcoach"); c != nil {
			out.Cheevos = append(out.Cheevos, *c)
		}
	}

	// Pass if no runs contain a fail.
	pass := !slices.ContainsFunc(runs, func(r hole.Run) bool { return !r.Pass })

	if pass && golfer != nil {
		out.RankUpdates =
			[]Golfer.RankUpdate{{Scoring: "bytes"}, {Scoring: "chars"}}

		if err := db.QueryRowContext(
			r.Context(),
			`SELECT earned,
			        failing_bytes,
			        old_bytes_joint, old_bytes_rank, old_bytes,
			        new_bytes_joint, new_bytes_rank, new_bytes,
			        new_bytes_solution_count,
			        old_best_bytes_first_golfer_id,
			        old_best_bytes_golfer_count,
			        old_best_bytes_golfer_id,
			        old_best_bytes,
			        old_best_bytes_submitted,
			        failing_chars,
			        old_chars_joint, old_chars_rank, old_chars,
			        new_chars_joint, new_chars_rank, new_chars,
			        new_chars_solution_count,
			        old_best_chars_first_golfer_id,
			        old_best_chars_golfer_count,
			        old_best_chars_golfer_id,
			        old_best_chars,
			        old_best_chars_submitted
			   FROM save_solution(
			            bytes   := CASE WHEN $3 = 'assembly'::lang
			                            THEN $5
			                            ELSE octet_length($1)
			                            END,
			            chars   := CASE WHEN $3 = 'assembly'::lang
			                            THEN NULL
			                            ELSE char_length($1)
			                            END,
			            code    := $1,
			            hole    := $2,
			            lang    := $3,
			            user_id := $4
			        )`,
			in.Code, in.Hole, in.Lang, golfer.ID, out.Runs[0].ASMBytes,
		).Scan(
			pq.Array(&out.Cheevos),
			&out.RankUpdates[0].FailingStrokes,
			&out.RankUpdates[0].From.Joint,
			&out.RankUpdates[0].From.Rank,
			&out.RankUpdates[0].From.Strokes,
			&out.RankUpdates[0].To.Joint,
			&out.RankUpdates[0].To.Rank,
			&out.RankUpdates[0].To.Strokes,
			&out.RankUpdates[0].NewSolutionCount,
			&out.RankUpdates[0].OldBestFirstGolferID,
			&out.RankUpdates[0].OldBestCurrentGolferCount,
			&out.RankUpdates[0].OldBestCurrentGolferID,
			&out.RankUpdates[0].OldBestStrokes,
			&out.RankUpdates[0].OldBestSubmitted,
			&out.RankUpdates[1].FailingStrokes,
			&out.RankUpdates[1].From.Joint,
			&out.RankUpdates[1].From.Rank,
			&out.RankUpdates[1].From.Strokes,
			&out.RankUpdates[1].To.Joint,
			&out.RankUpdates[1].To.Rank,
			&out.RankUpdates[1].To.Strokes,
			&out.RankUpdates[1].NewSolutionCount,
			&out.RankUpdates[1].OldBestFirstGolferID,
			&out.RankUpdates[1].OldBestCurrentGolferCount,
			&out.RankUpdates[1].OldBestCurrentGolferID,
			&out.RankUpdates[1].OldBestStrokes,
			&out.RankUpdates[1].OldBestSubmitted,
		); err != nil {
			panic(err)
		}

		recordUpdates := make([]Golfer.RankUpdate, 0, 2)

		for _, rank := range out.RankUpdates {
			if rank.From.Strokes.V == rank.To.Strokes.V {
				continue
			}

			// This keeps track of which updates (if any) represent new records or diamond matches.
			if rank.To.Rank.V == 1 {
				if !rank.To.Joint.V ||
					rank.OldBestCurrentGolferCount.Valid && rank.OldBestCurrentGolferCount.V == 1 {
					recordUpdates = append(recordUpdates, rank)
				}
			}
		}

		// If any of the updates are record breakers, announce them on Discord
		if len(recordUpdates) > 0 {
			go discord.LogNewRecord(golfer, holeObj, langObj, recordUpdates, db)
		}

		// For now don't show any popups for experimental solutions.
		if experimental {
			out.RankUpdates = []Golfer.RankUpdate{}
		}

		// Cheevos.
		if experimental {
			if c := golfer.Earn(db, "black-box-testing"); c != nil {
				out.Cheevos = append(out.Cheevos, *c)
			}

			if experimentalHole && experimentalLang {
				if c := golfer.Earn(db, "double-slit-experiment"); c != nil {
					out.Cheevos = append(out.Cheevos, *c)
				}
			}
		} else {
			// TODO Use the golfer's timezone from /settings.
			// TODO Move these to save_solution() in the DB.
			var (
				now   = time.Now().UTC()
				month = now.Month()
				day   = now.Day()
			)

			if month == time.October && day == 2 {
				if c := golfer.Earn(db, "happy-birthday-code-golf"); c != nil {
					out.Cheevos = append(out.Cheevos, *c)
				}
			}

			switch in.Hole {
			case "12-days-of-christmas":
				if (month == time.December && day >= 25) ||
					(month == time.January && day <= 5) {
					if c := golfer.Earn(db, "twelvetide"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			case "star-wars-gpt", "star-wars-opening-crawl":
				if month == time.May && day == 4 {
					if c := golfer.Earn(db, "may-the-4ᵗʰ-be-with-you"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			case "tic-tac-toe":
				if month == time.February && day == 14 {
					if c := golfer.Earn(db, "hugs-and-kisses"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			case "united-states":
				if month == time.July && day == 4 {
					if c := golfer.Earn(db, "independence-day"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			case "vampire-numbers":
				if month == time.October && day == 31 {
					if c := golfer.Earn(db, "vampire-byte"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			case "π":
				if month == time.March && day == 14 {
					if c := golfer.Earn(db, "pi-day"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			case "τ":
				if month == time.June && day == 28 {
					if c := golfer.Earn(db, "how-about-second-pi"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			}

			if in.Lang == "viml" && golfer.Settings["hole"]["editor-keymap"] == "vim" {
				if c := golfer.Earn(db, "real-programmers"); c != nil {
					out.Cheevos = append(out.Cheevos, *c)
				}
			}

			if !golfer.Earned("smörgåsbord") {
				var earn bool
				if err := db.Get(
					&earn,
					`WITH distinct_holes AS (
				    SELECT DISTINCT hole
				      FROM solutions
				     WHERE NOT failing AND user_id = $1
				) SELECT (
				    SELECT COUNT(DISTINCT $2::hstore->hole::text)
				      FROM distinct_holes
				) = (
				    SELECT COUNT(DISTINCT cat)
				      FROM unnest(avals($2)) cat
				)`,
					golfer.ID,
					config.HoleCategoryHstore,
				); err != nil {
					panic(err)
				}

				if earn {
					if c := golfer.Earn(db, "smörgåsbord"); c != nil {
						out.Cheevos = append(out.Cheevos, *c)
					}
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	// Avoid returning "null" arrays. Eventually fix with json/v2.
	// See https://github.com/golang/go/issues/37711#issuecomment-1750018405
	if out.Cheevos == nil {
		out.Cheevos = []config.Cheevo{}
	}
	if out.RankUpdates == nil {
		out.RankUpdates = []Golfer.RankUpdate{}
	}

	if err := enc.Encode(&out); err != nil {
		panic(err)
	}
}

// GET /mini-rankings/{hole}/{lang}/{scoring:bytes|chars}/{view:top|me|following}
func apiMiniRankingsGET(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if r.FormValue("long") == "1" {
		limit = 99
	}

	var (
		holeID  = param(r, "hole")
		langID  = param(r, "lang")
		scoring = param(r, "scoring")
		view    = param(r, "view")
	)

	// No need to check scoring & view, they're covered by chi.
	hole := config.AllHoleByID[holeID]
	lang := config.AllLangByID[langID]
	if hole == nil || lang == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	otherScoring := "bytes"
	if scoring == "bytes" {
		otherScoring = "chars"
	}

	var followLimit, userID int
	if golfer := session.Golfer(r); golfer != nil {
		followLimit = golfer.FollowLimit()
		userID = golfer.ID
	}

	type entry struct {
		Bytes      *int `json:"bytes"`
		BytesChars *int `json:"bytes_chars"`
		Chars      *int `json:"chars"`
		CharsBytes *int `json:"chars_bytes"`
		Golfer     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"golfer"`
		Me   bool `json:"me"`
		Rank int  `json:"rank"`
	}

	sqlWhere, sqlLimit := "true", "$6"
	switch view {
	case "me":
		sqlWhere = "row > COALESCE((SELECT row FROM ranks WHERE me), 0) - $6"
		sqlLimit = "$6 * 2"
	case "following":
		sqlWhere = fmt.Sprint("user_id = ANY(following($1, ", followLimit, "))")
	}

	// We don't use the rankings view as we want instant updates upon solution
	// submit, therefore we skip scoring to keep it fast.
	entries := make([]entry, 0, limit)
	if err := session.Database(r).Select(
		&entries,
		`WITH ranks AS (
		    SELECT ROW_NUMBER() OVER (ORDER BY `+scoring+`, submitted) row,
		           RANK()       OVER (ORDER BY `+scoring+`),
		           user_id,
		           `+scoring+`,
		           `+otherScoring+` `+scoring+`_`+otherScoring+`,
		           user_id = $1 me
		      FROM solutions
		     WHERE hole = $2
		       AND lang = $3
		       AND scoring = $4
		       AND NOT failing
		), other_scoring AS (
		    SELECT user_id,
		           `+otherScoring+`,
		           `+scoring+` `+otherScoring+`_`+scoring+`
		      FROM solutions
		     WHERE hole = $2
		       AND lang = $3
		       AND scoring = $5
		       AND NOT failing
		)   SELECT bytes, bytes_chars, chars, chars_bytes, me, rank,
		           id "golfer.id", login "golfer.name"
		      FROM ranks
		      JOIN users ON id = user_id
		 LEFT JOIN other_scoring USING(user_id)
		     WHERE `+sqlWhere+`
		  ORDER BY row
		     LIMIT `+sqlLimit,
		userID,
		holeID,
		langID,
		scoring,
		otherScoring,
		limit,
	); err != nil {
		panic(err)
	}

	// Trim the rows to limit with "me" as close to the middle as possible.
	// TODO It would simplify everything if we could fold this into the SQL.
	length := len(entries)
	if view == "me" && length > limit {
		me := slices.IndexFunc(entries, func(e entry) bool { return e.Me })
		// Before: me entries, then "me" entry, then len(entries)-me-1 entries
		// 	with me <= limit; len(entries) <= 2*limit; len(entries)-me-1 <= limit-1
		if me <= limit/2 {
			entries = entries[:limit]
		} else if me >= length-(limit+1)/2 {
			// Impossible case?
			entries = entries[length-limit:]
		} else {
			entries = entries[me-limit/2 : me+(limit+1)/2]
		}
	}

	encodeJSON(w, entries)
}
