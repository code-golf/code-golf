package routes

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/code-golf/code-golf/session"
	"github.com/code-golf/code-golf/trophy"
)

type progress struct{ Completed, Steps int }

func (p progress) Percent() int { return p.Completed * 100 / p.Steps }

func (p progress) SVG() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<svg viewbox="0 0 8 8"><circle cx="4" cy="4" r="3.5" stroke-dasharray="22 22" stroke-dashoffset="%d"/></svg>`,
		22-p.Completed*22/p.Steps,
	))
}

// GolferAchievements serves GET /golfers/{golfer}/achievements
func GolferAchievements(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Cheevos  map[string][]*trophy.Trophy
		Progress map[string]*progress
	}{
		trophy.Tree,
		map[string]*progress{
			"caffeinated":              {0, 2},
			"elephpant-in-the-room":    {0, 1},
			"happy-birthday-code-golf": {0, 1},
			"hello-world":              {0, 1},
			"inception":                {0, 1},
			"independence-day":         {0, 1},
			"interview-ready":          {0, 1},
			"its-over-9000":            {0, 9001},
			"my-god-its-full-of-stars": {0, 1},
			"ouroboros":                {0, 1},
			"patches-welcome":          {0, 1},
			"pi-day":                   {0, 1},
			"polyglot":                 {0, 10},
			"slowcoach":                {0, 1},
			"the-watering-hole":        {0, 19},
			"tim-toady":                {0, 2},
			"twelvetide":               {0, 1},
		},
	}

	// Random progress.
	for _, cheevo := range data.Progress {
		cheevo.Completed = rand.Intn(cheevo.Steps + 1)
	}

	render(w, r, "golfer/achievements", session.GolferInfo(r).Golfer.Name, data)
}
