package routes

import (
	"html/template"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pretty"
)

var nextHole = config.ExpHoleByID["minesweeper"]

type banner struct {
	Body          template.HTML
	HideKey, Type string
}

func banners(golfer *golfer.Golfer, now time.Time) (banners []banner) {
	// Upcoming hole.
	if hole := nextHole; hole != nil {
		t := hole.Released.AsTime(time.UTC)
		if golfer != nil {
			t = t.In(golfer.Location())
		}

		in := "in approximately " + pretty.Time(t)
		if strings.Contains(string(in), "ago") {
			in = "momentarily"
		}

		banners = append(banners, banner{
			HideKey: "upcoming-hole-" + hole.ID,
			Type:    "info",
			Body: "The <a href=/" + template.HTML(hole.ID) + ">" +
				template.HTML(hole.Name) + "</a> hole will go live " + in +
				". Why not try and solve it ahead of time?",
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
	delta := now.Year() - 2000

	location := golfer.Location()

Cheevo:
	for _, cheevo := range config.CheevoList {
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

				break Cheevo
			}
		}
	}

	// Latest hole (if unsolved).
	if hole := config.RecentHoles[0]; !golfer.Solved(hole.ID) {
		banners = append(banners, banner{
			HideKey: "latest-hole-" + hole.ID,
			Type:    "info",
			Body: template.HTML(`The <a href="/` + hole.ID + `">` +
				hole.Name + "</a> hole is now live! Why not try and solve it?"),
		})
	}

	return
}
