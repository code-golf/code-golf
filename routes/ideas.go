package routes

import (
	"cmp"
	"net/http"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

// GET /ideas
func ideasGET(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		Category, Title          string
		Hole                     *config.Hole
		ID, ThumbsDown, ThumbsUp int
		Lang                     *config.Lang
	}

	var data struct{ Holes, Ideas, Langs []idea }

	// FIXME This is only a slice because of semiprime/sphenic.
	holes := map[int][]*config.Hole{}
	for _, hole := range config.ExpHoleList {
		holes[hole.Experiment] = append(holes[hole.Experiment], hole)
	}

	langs := map[int]*config.Lang{}
	for _, lang := range config.ExpLangList {
		langs[lang.Experiment] = lang
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
				i.Hole = hole
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
