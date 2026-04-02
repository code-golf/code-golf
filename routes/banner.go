package routes

import (
	"bytes"
	"html/template"
	"maps"
	"slices"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pretty"
	"github.com/code-golf/code-golf/views"
)

type banner struct {
	Body          template.HTML
	HideKey, Type string
}

func banners(golfer *golfer.Golfer) (banners []banner) {
	// Upcoming holes.
	if holes := config.NextHoles; len(holes) != 0 {
		t := holes[0].Released.AsTime(time.UTC)
		if golfer != nil {
			t = t.In(golfer.Location())
		}

		var html bytes.Buffer
		if err := views.Render(&html, "banners/upcoming-holes", struct {
			Holes []*config.Hole
			Time  time.Time
		}{holes, t}); err != nil {
			panic(err)
		}

		banners = append(banners, banner{
			Body:    template.HTML(html.String()),
			HideKey: "upcoming-hole-" + holes[0].ID,
			Type:    "info",
		})
	}

	// Currently all the global banners require a golfer.
	if golfer == nil {
		return
	}

	// Pending account deletion.
	if delete := golfer.Delete; delete.Valid {
		banners = append(banners, banner{
			Type: "alert",
			Body: template.HTML(
				"Your account will be permanently deleted on the " +
					delete.V.Format("2 Jan 2006") + "." +
					"<p>If you wish to stop this, visit " +
					"<a href=/golfer/settings/delete-account>settings</a> " +
					"and cancel the deletion."),
		})
	}

	// Failing solutions.
	if failing := golfer.FailingSolutions; len(failing) > 0 {
		banner := banner{
			Type: "alert",
			Body: "The following solutions of yours have been marked as " +
				"failing and no longer contribute to scoring; " +
				"please update them to pass:<ul>",
		}

		type GroupItem struct {
			HoleID, LangID, KeyName, OtherName string
		}

		byLang := make(map[string][]GroupItem)
		byHole := make(map[string][]GroupItem)

		for _, solution := range failing {
			langName := config.AllLangByID[solution.Lang].Name
			holeName := config.AllHoleByID[solution.Hole].Name
			byLang[solution.Lang] = append(byLang[solution.Lang], GroupItem{solution.Hole, solution.Lang, langName, holeName})
			byHole[solution.Hole] = append(byHole[solution.Hole], GroupItem{solution.Hole, solution.Lang, holeName, langName})
		}

		groups := byHole
		if len(byLang) < len(byHole) {
			groups = byLang
		}
		keys := slices.Collect(maps.Keys(groups))
		slices.Sort(keys)

		for _, key := range keys {
			solutionGroup := groups[key]
			for index, solution := range solutionGroup {
				if index > 0 {
					banner.Body += template.HTML(", ")
				} else {
					banner.Body += template.HTML("<li>" + solution.KeyName + ": ")
				}
				banner.Body += template.HTML(`<a href="/` + solution.HoleID + "#" + solution.LangID + `">` + solution.OtherName + "</a>")
			}
		}

		banner.Body += "</ul>"

		banners = append(banners, banner)
	}

	// Our date-specific cheevos are set around the year 2000.
	now := time.Now().UTC()
	delta := now.Year() - 2000

	location := golfer.Location()

	for _, cheevo := range config.CheevoTree["Date Specific"] {
		if golfer.Earned(cheevo.ID) {
			continue
		}

		for i := 0; i < len(cheevo.Times); i += 2 {
			start := cheevo.Times[i].AddDate(delta, 0, 0)
			end := cheevo.Times[i+1].AddDate(delta, 0, 0)

			var body template.HTML
			var hideKey string
			if start.AddDate(0, 0, -7).Before(now) && now.Before(start) {
				body = "be available in " +
					pretty.Time(start.In(location)) + "."
				hideKey = "cheevo-before-" + start.Format(time.DateOnly) + "-" + cheevo.ID
			} else if start.Before(now) && now.Before(end) {
				body = "stop being available in " +
					pretty.Time(end.In(location)) + "."
				hideKey = "cheevo-until-" + end.Format(time.DateOnly) + "-" + cheevo.ID
			}

			if body != "" {
				body = template.HTML("The "+cheevo.Emoji+" <b>"+
					cheevo.Name+"</b> achievement will ") + body +
					"<p>" + cheevo.Description

				banners = append(banners, banner{
					Body:    body,
					HideKey: hideKey,
					Type:    "info",
				})

				break
			}
		}
	}

	// Latest hole(s) (if unsolved).
	if !golfer.SolvedLatestHole {
		var html bytes.Buffer
		if err := views.Render(
			&html, "banners/latest-holes", config.LatestHoles,
		); err != nil {
			panic(err)
		}

		banners = append(banners, banner{
			Body:    template.HTML(html.String()),
			HideKey: "latest-hole-" + config.LatestHoles[0].ID,
			Type:    "info",
		})
	}

	// Latest lang (if unsolved and still fresh).
	if lang := config.LatestLang; !golfer.SolvedLatestLang && lang.Fresh() {
		banners = append(banners, banner{
			HideKey: "latest-lang-" + lang.ID,
			Type:    "info",
			Body: template.HTML(lang.Name +
				" is now live! Why not try and solve a hole in it?"),
		})
	}

	// Hole of the Week.
	if html, week := config.HoleOfTheWeek(); html != "" && !golfer.SolvedHoleOfTheWeek {
		banners = append(banners, banner{
			HideKey: "hole-of-the-week-" + week.Format(time.DateOnly),
			Type:    "info",
			Body:    html,
		})
	}

	return
}
