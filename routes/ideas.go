package routes

import (
	"cmp"
	"net/http"
	"slices"
	"strings"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/session"
)

func ideaColor(kind string) string {
	switch kind {
	case "hole":
		return "purple"
	case "lang":
		return "yellow"
	case "cheevo":
		return "green"
	}
	return "red"
}

// GET /ideas
func ideasGET(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		Hole                     *config.Hole
		ID, ThumbsDown, ThumbsUp int
		Title, Kind, KindColor   string
	}

	data := struct {
		Holes, Ideas []idea
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
		}

		i.KindColor = ideaColor(i.Kind)
		data.Ideas = append(data.Ideas, i)
	}

	slices.SortStableFunc(data.Holes, func(a, b idea) int {
		return cmp.Compare(b.ThumbsUp-b.ThumbsDown, a.ThumbsUp-a.ThumbsDown)
	})

	render(w, r, "ideas", data, "Ideas")
}
