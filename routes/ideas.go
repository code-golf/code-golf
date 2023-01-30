package routes

import (
	"net/http"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
	"golang.org/x/exp/slices"
)

// GET /ideas
func ideasGET(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		Hole                     *config.Hole
		ID, ThumbsDown, ThumbsUp int
		Title                    string
	}

	data := struct {
		Holes, Ideas, Langs []idea
	}{Holes: make([]idea, len(config.ExpHoleList))}

	for i, hole := range config.ExpHoleList {
		data.Holes[i] = idea{Hole: hole, ID: hole.Experiment, Title: hole.Name}
	}

	var ideas []idea
	if err := session.Database(r).Select(
		&ideas,
		"SELECT * FROM ideas ORDER BY thumbs_up - thumbs_down DESC, title",
	); err != nil {
		panic(err)
	}

rows:
	for _, i := range ideas {
		for j, hole := range data.Holes {
			if hole.ID == i.ID {
				data.Holes[j].ThumbsDown = i.ThumbsDown
				data.Holes[j].ThumbsUp = i.ThumbsUp
				continue rows
			}
		}

		if strings.HasPrefix(i.Title, "Add ") && strings.HasSuffix(i.Title, " Lang") {
			i.Title = i.Title[len("Add ") : len(i.Title)-len(" Lang")]

			data.Langs = append(data.Langs, i)
		} else {
			data.Ideas = append(data.Ideas, i)
		}
	}

	slices.SortStableFunc(data.Holes, func(a, b idea) bool {
		return a.ThumbsUp-a.ThumbsDown > b.ThumbsUp-b.ThumbsDown
	})

	render(w, r, "ideas", data, "Ideas")
}
