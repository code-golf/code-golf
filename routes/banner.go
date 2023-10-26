package routes

import (
	"html/template"
	"strings"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pretty"
)

type banner struct {
	Body template.HTML
	Type string
}

// TODO Allow a golfer to hide individual banners #709.
func banners(golfer *golfer.Golfer, now time.Time) (banners []banner) {
	in := "in " + pretty.Time(time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC))
	if strings.Contains(string(in), "ago") {
		in = "momentarily"
	}

	hole := config.ExpHoleByID["dfa-simulator"]
	banners = append(banners, banner{
		Type: "info",
		Body: "The <a href=/" + template.HTML(hole.ID) + ">" +
			template.HTML(hole.Name) + "</a> hole will go live " + in +
			". Why not try and solve it ahead of time?",
	})

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
					delete.Time.Format("2 Jan 2006") + "." +
					"<p>If you wish to stop this, visit " +
					"<a href=/golfer/settings>settings</a> " +
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

		prevHoleID := ""
		for _, solution := range failing {
			hole := config.HoleByID[solution.Hole]
			lang := config.LangByID[solution.Lang]

			if prevHoleID == hole.ID {
				banner.Body += template.HTML(`, <a href="/` + hole.ID + "#" +
					lang.ID + `">` + lang.Name + "</a>")
			} else {
				banner.Body += template.HTML("<li>" + hole.Name + `: <a href="/` + hole.ID + "#" +
					lang.ID + `">` + lang.Name + "</a>")
			}
			prevHoleID = hole.ID
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
			if start.AddDate(0, 0, -7).Before(now) && now.Before(start) {
				body = "be available in " +
					pretty.Time(start.In(location)) + "."
			} else if start.Before(now) && now.Before(end) {
				body = "stop being available in " +
					pretty.Time(end.In(location)) + "."
			}

			if body != "" {
				body = template.HTML("The "+cheevo.Emoji+" <b>"+
					cheevo.Name+"</b> achievement will ") + body

				banners = append(banners, banner{body, "info"})

				break Cheevo
			}
		}
	}

	// Latest hole (if unsolved).
	if hole := config.RecentHoles[0]; !golfer.Solved(hole.ID) {
		banners = append(banners, banner{
			Type: "info",
			Body: template.HTML(`The <a href="/` + hole.ID + `">` +
				hole.Name + "</a> hole is now live! Why not try and solve it?"),
		})
	}

	return
}
