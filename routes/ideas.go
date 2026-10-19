package routes

import (
	"cmp"
	"net/http"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /ideas
func ideasGET(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		Category, Title          string
		Solved                   bool
		Hole                     *config.Hole
		Golfers                  int
		ID, ThumbsDown, ThumbsUp int
		Lang                     *config.Lang
	}

	var data struct{ Holes, Ideas, Langs []idea }

	// FIXME This is only a slice because of semiprime/sphenic.
	holes := map[int][]*config.Hole{}
	holeIDs := make([]string, len(config.ExpHoleList))
	for i, hole := range config.ExpHoleList {
		holes[hole.Experiment] = append(holes[hole.Experiment], hole)
		holeIDs[i] = hole.ID
	}

	// How many golfers have a passing solution for each experimental hole.
	var solves []struct {
		Golfers int
		Hole    string
	}
	if err := session.Database(r).Select(
		&solves,
		`SELECT COUNT(DISTINCT user_id) golfers, hole
		   FROM solutions
		  WHERE hole = ANY($1) AND NOT failing
	   GROUP BY hole`,
		pq.Array(holeIDs),
	); err != nil {
		panic(err)
	}

	golfers := map[string]int{}
	for _, solve := range solves {
		golfers[solve.Hole] = solve.Golfers
	}

	langs := map[int]*config.Lang{}
	for _, lang := range config.ExpLangList {
		langs[lang.Experiment] = lang
	}

	// Experimental holes the golfer has a passing solution for.
	solved := map[string]bool{}
	if golfer := session.Golfer(r); golfer != nil {
		var holes []string
		if err := session.Database(r).Select(
			&holes,
			"SELECT DISTINCT hole FROM solutions WHERE user_id = $1 AND NOT failing",
			golfer.ID,
		); err != nil {
			panic(err)
		}

		for _, hole := range holes {
			solved[hole] = true
		}
	}

	var ideas []idea
	if err := session.Database(r).Select(
		&ideas,
		"SELECT * FROM ideas ORDER BY thumbs_up - thumbs_down DESC, title",
	); err != nil {
		panic(err)
	}

	for _, i := range ideas {
		if holes, ok := holes[i.ID]; ok {
			for _, hole := range holes {
				i.Golfers = golfers[hole.ID]
				i.Hole = hole
				i.Solved = solved[hole.ID]
				i.Title = hole.Name
				data.Holes = append(data.Holes, i)
			}
		} else if lang, ok := langs[i.ID]; ok {
			i.Lang = lang
			i.Title = lang.Name
			data.Langs = append(data.Langs, i)
		} else {
			data.Ideas = append(data.Ideas, i)
		}
	}

	// Sort by vote difference, then case-insensitive title.
	for _, ideas := range [][]idea{data.Holes, data.Langs} {
		slices.SortFunc(ideas, func(a, b idea) int {
			return cmp.Or(
				cmp.Compare(b.ThumbsUp-b.ThumbsDown, a.ThumbsUp-a.ThumbsDown),
				cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title)),
			)
		})
	}

	render(w, r, "ideas", data, "Ideas")
}
