package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}/curves
func golferCurvesGET(w http.ResponseWriter, r *http.Request) {
	type improvement struct {
		Bytes     int       `db:"bytes"     json:"bytes"`
		Chars     *int      `db:"chars"     json:"chars"`
		Hole      string    `db:"hole"      json:"hole"`
		Lang      string    `db:"lang"      json:"lang"`
		Scoring   string    `db:"scoring"   json:"scoring"`
		Submitted time.Time `db:"submitted" json:"submitted"`
	}

	// Sb of the Bayesian estimator in the rankings, so the client can work out
	// how many points a past solution would be worth on today's leaderboards.
	type estimator struct {
		Hole    string  `db:"hole"    json:"hole"`
		Lang    string  `db:"lang"    json:"lang"`
		Scoring string  `db:"scoring" json:"scoring"`
		Sb      float64 `db:"sb"      json:"sb"`
	}

	type medal struct {
		Lang      string    `db:"lang"      json:"lang"`
		Medal     string    `db:"medal"     json:"medal"`
		Scoring   string    `db:"scoring"   json:"scoring"`
		Submitted time.Time `db:"submitted" json:"submitted"`
	}

	data := struct {
		Cheevos      []time.Time       `json:"cheevos"`
		Estimators   []estimator       `json:"estimators"`
		Improvements []improvement     `json:"improvements"`
		LangNames    map[string]string `json:"langNames"`
		Medals       []medal           `json:"medals"`
	}{LangNames: map[string]string{}}

	db := session.Database(r)
	golferID := session.GolferInfo(r).ID

	// Every improvement of the golfer's current solutions. Solutions older than
	// the log only have their current submission.
	if err := db.Select(
		&data.Improvements,
		`WITH solutions AS (
		    SELECT hole, lang, scoring, submitted, bytes, chars
		      FROM stable_passing_solutions
		     WHERE user_id = $1
		)   SELECT log.submitted, hole, lang, scoring, log.bytes, log.chars
		      FROM solutions_log log
		      JOIN solutions USING(hole, lang, scoring)
		     WHERE log.user_id = $1
		 UNION ALL
		    SELECT submitted, hole, lang, scoring, bytes, chars
		      FROM solutions
		     WHERE NOT EXISTS (
		        SELECT FROM solutions_log log
		         WHERE log.user_id = $1
		           AND (log.hole, log.lang, log.scoring)
		             = (solutions.hole, solutions.lang, solutions.scoring)
		     )
		  ORDER BY submitted`,
		golferID,
	); err != nil {
		panic(err)
	}

	for _, improvement := range data.Improvements {
		data.LangNames[improvement.Lang] = config.AllLangByID[improvement.Lang].Name
	}

	if err := db.Select(
		&data.Estimators,
		`SELECT hole, lang, scoring, points * strokes / 1000.0 sb
		   FROM rankings
		  WHERE user_id = $1 AND NOT experimental`,
		golferID,
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Medals,
		`SELECT lang, medal, scoring, submitted
		   FROM medals
		  WHERE user_id = $1`,
		golferID,
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Cheevos,
		"SELECT earned FROM cheevos WHERE user_id = $1",
		golferID,
	); err != nil {
		panic(err)
	}

	render(w, r, "golfer/curves", data, session.GolferInfo(r).Name)
}
