package routes

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/code-golf/code-golf/session"
)

type cheevo struct {
	Emoji, Title string
	Description  template.HTML
	Value, Max   int
}

func (c cheevo) Percent() int { return c.Value * 100 / c.Max }

func (c cheevo) Progress() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<svg viewbox="0 0 8 8"><circle cx="4" cy="4" r="3.5" stroke-dasharray="22 22" stroke-dashoffset="%d"/></svg>`,
		22-c.Value*22/c.Max,
	))
}

// GolferAchievements serves GET /golfers/{golfer}/achievements
func GolferAchievements(w http.ResponseWriter, r *http.Request) {
	cheevos := []struct {
		Category string
		Cheevos  []cheevo
		Earned   int
	}{
		{
			Category: "Progression",
			Cheevos: []cheevo{
				{
					"👋", "Hello, World!",
					"Solve your first hole.",
					0, 1,
				},
				{
					"🍞", "Baker’s Dozen",
					"Solve your thirteenth hole.",
					0, 13,
				},
				{
					"🍺", "The Watering Hole",
					"Solve your nineteenth hole.",
					0, 19,
				},
				{
					"👍", "DON’T PANIC!",
					"Solve your forty-second hole.",
					0, 42,
				},
				{
					"🐉", "It’s Over 9000!",
					"Earn over 9,000 points.",
					0, 9001,
				},
			},
		},
		{
			Category: "Hole/Lang Specific",
			Cheevos: []cheevo{
				{
					"💼", "Interview Ready",
					"Solve <a href=/fizz-buzz>Fizz Buzz</a>.",
					0, 1,
				},
				{
					"🐘", "ElePHPant in the Room",
					"Solve any hole in PHP.",
					0, 1,
				},
				{
					"🐍", "Ouroboros",
					"Solve <a href=/quine#python>Quine in Python</a>.",
					0, 1,
				},
				{
					"🐪", "Tim Toady",
					"Solve any hole in both Perl and Raku.",
					0, 2,
				},
				{
					"☕", "Caffeinated",
					"Solve any hole in both Java and JavaScript.",
					0, 2,
				},
				{
					"🧠", "Inception",
					"Solve <a href=/brainfuck#brainfuck>Brainfuck in Brainfuck</a>.",
					0, 1,
				},
			},
		},
		{
			Category: "Date Specific",
			Cheevos: []cheevo{
				{
					"🎂", "Happy Birthday, Code Golf",
					"Solve any hole on " +
						"<a href=//github.com/code-golf/code-golf/commit/4b44>" +
						"2 Oct</a>.",
					0, 1,
				},
				{
					"🇺🇸", "Independence Day",
					"Solve <a href=/united-states>United States</a> on " +
						"<a href=//www.wikipedia.org/wiki/Independence_Day_(United_States)>" +
						"4 Jul</a>.",
					0, 1,
				},
				{
					"🥧", "Pi Day",
					"Solve <a href=/π>π</a> on " +
						"<a href=//www.wikipedia.org/wiki/Pi_Day>14 Mar.</a>",
					0, 1,
				},
				{
					"🎅", "Twelvetide",
					"Solve <a href=/12-days-of-christmas>12 Days of Christmas</a> during <a href=//www.wikipedia.org/wiki/Twelvetide>25 Dec – 5 Jan</a>.",
					0, 1,
				},
			},
		},
		{
			Category: "Miscellaneous",
			Cheevos: []cheevo{
				{
					"🦥", "Slowcoach",
					"Fail any hole by exceeding the time limit.",
					0, 1,
				},
				{
					"📕", "RTFM",
					"Vist the <a href=/about>About page</a>.",
					0, 1,
				},
			},
		},
		{
			Category: "GitHub",
			Cheevos: []cheevo{
				{
					"⭐", "My God, It’s Full of Stars",
					"Star <a href=//github.com/code-golf/code-golf> the " +
						"Code Golf repository</a>.",
					0, 1,
				},
				{
					"💾", "Patches Welcome",
					"Contribute a merged PR to " +
						"<a href=//github.com/code-golf/code-golf> the " +
						"Code Golf repository </a>.",
					0, 1,
				},
			},
		},
	}

	// Random earned values.
	for i, category := range cheevos {
		for j, cheevo := range category.Cheevos {
			category.Cheevos[j].Value = rand.Intn(cheevo.Max + 1)

			if category.Cheevos[j].Value == cheevo.Max {
				cheevos[i].Earned++
			}
		}
	}

	render(w, r, "golfer/achievements", session.GolferInfo(r).Golfer.Name, cheevos)
}
