package routes

import (
	"database/sql"
	json_v1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/code-golf/code-golf/config"
	Golfer "github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
	"golang.org/x/sync/semaphore"
)

// A polyglot runs in every standard lang, limit how many langs run at once so
// that polyglots don't starve regular solutions of the shared run semaphore.
var polyglotSemaphore = semaphore.NewWeighted(3)

// polyglotEvent is one line of the NDJSON stream sent back by POST /polyglot.
// The first line has Total, then each lang gets a "compile" line and, if it
// compiled, a "test" line. The last line has Done.
type polyglotEvent struct {
	Total int `json:"total,omitzero"`

	Lang  string     `json:"lang,omitzero"`
	Phase string     `json:"phase,omitzero"`
	OK    bool       `json:"ok,omitzero"`
	Runs  []hole.Run `json:"runs,omitzero"`

	Done     bool     `json:"done,omitzero"`
	Langs    []string `json:"langs,omitzero"`
	LoggedIn bool     `json:"logged_in,omitzero"`
	Saved    bool     `json:"saved,omitzero"`
}

// POST /polyglot
func polyglotPOST(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Code string `json:"code"`
		Hole string `json:"hole"`
	}

	if err := json.UnmarshalRead(r.Body, &in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	holeObj := config.AllHoleByID[in.Hole]
	if holeObj == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 128 KiB, >= because arguments needs a null termination.
	if len(in.Code) >= 128*1024 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	golfer := session.Golfer(r)

	var mux sync.Mutex
	send := func(e polyglotEvent) {
		b, err := json.Marshal(
			e, jsontext.AllowInvalidUTF8(true), json_v1.FormatDurationAsNano(true))
		if err != nil {
			panic(err)
		}

		mux.Lock()
		defer mux.Unlock()

		w.Write(append(b, '\n'))
		w.(http.Flusher).Flush()
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	send(polyglotEvent{Total: len(config.LangList)})

	// Compile (run the first test case), then test the langs that compiled.
	passed := make([]bool, len(config.LangList))
	var wg sync.WaitGroup
	for i, lang := range config.LangList {
		wg.Go(func() {
			if err := polyglotSemaphore.Acquire(r.Context(), 1); err != nil {
				return
			}
			defer polyglotSemaphore.Release(1)

			run, err := hole.PlayFirst(r.Context(), holeObj, lang, in.Code)
			if err != nil {
				log.Printf("polyglot %s/%s: %v", holeObj.ID, lang.ID, err)
				return
			}

			compiled := !run.Timeout && run.ExitCode == 0
			send(polyglotEvent{
				Lang: lang.ID, Phase: "compile", OK: compiled, Runs: []hole.Run{run},
			})
			if !compiled {
				return
			}

			runs, err := hole.Play(r.Context(), holeObj, lang, in.Code)
			if err != nil {
				log.Printf("polyglot %s/%s: %v", holeObj.ID, lang.ID, err)
				return
			}

			passed[i] = true
			for _, run := range runs {
				passed[i] = passed[i] && run.Pass
			}

			send(polyglotEvent{
				Lang: lang.ID, Phase: "test", OK: passed[i], Runs: runs,
			})
		})
	}
	wg.Wait()

	// The client went away, there's no one to tell and nothing to save.
	if r.Context().Err() != nil {
		return
	}

	done := polyglotEvent{Done: true, Langs: []string{}, LoggedIn: golfer != nil}
	for i, lang := range config.LangList {
		if passed[i] {
			done.Langs = append(done.Langs, lang.ID)
		}
	}

	// Keep the golfer's best polyglot: more langs, then fewer bytes.
	if golfer != nil && len(done.Langs) > 0 && holeObj.Experiment == 0 {
		if err := session.Database(r).QueryRow(
			`INSERT INTO polyglot_solutions (user_id, hole, code, langs)
			      VALUES ($1, $2, $3, $4::lang[])
			 ON CONFLICT (user_id, hole) DO UPDATE
			         SET code = excluded.code,
			             langs = excluded.langs,
			             submitted = excluded.submitted
			       WHERE cardinality(excluded.langs) > cardinality(polyglot_solutions.langs)
			          OR cardinality(excluded.langs) = cardinality(polyglot_solutions.langs)
			         AND octet_length(excluded.code) < octet_length(polyglot_solutions.code)
			   RETURNING true`,
			golfer.ID, holeObj.ID, in.Code, pq.Array(done.Langs),
		).Scan(&done.Saved); err != nil && !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}
	}

	send(done)
}

// GET /api/polyglot-rankings/{hole}
func apiPolyglotRankingsGET(w http.ResponseWriter, r *http.Request) {
	holeObj := config.AllHoleByID[param(r, "hole")]
	if holeObj == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var userID int
	if golfer := session.Golfer(r); golfer != nil {
		userID = golfer.ID
	}

	type entry struct {
		Bytes  int               `json:"bytes"`
		Golfer Golfer.GolferLink `json:"golfer"`
		Langs  int               `json:"langs"`
		Me     bool              `json:"me"`
		Rank   int               `json:"rank"`
	}

	entries := []entry{}
	if err := session.Database(r).Select(
		&entries,
		` SELECT RANK() OVER (ORDER BY cardinality(langs) DESC, octet_length(code)) rank,
		         cardinality(langs) langs,
		         octet_length(code) bytes,
		         user_id = $1       me,
		         avatar_url         "golfer.avatar_url",
		         name               "golfer.name"
		    FROM polyglot_solutions
		    JOIN golfers_with_avatars ON id = user_id
		   WHERE hole = $2
		ORDER BY rank, submitted
		   LIMIT 99`,
		userID, holeObj.ID,
	); err != nil {
		panic(err)
	}

	encodeJSON(w, entries)
}
