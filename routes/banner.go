package routes

import (
	"html/template"
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
func banners(golfer *golfer.Golfer) (banners []banner) {
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

		prevHole := "none"
		for _, solution := range failing {
			hole := config.HoleByID[solution.Hole]
			lang := config.LangByID[solution.Lang]

			if prevHole == hole {
				banner.Body += template.HTML(`, <a href="/` + hole.ID + "#" +
					lang.ID + `">` + lang.Name + "</a>")
			} else {
				banner.Body += template.HTML("<li>" + lang.Name + `: <a href="/` + hole.ID + "#" +
					lang.ID + `">` + lang.Name + "</a>")
			}
			prevHole = hole
		}

		banner.Body += "</ul>"

		banners = append(banners, banner)
	}

	// Current cheevo. TODO Generalise.
	if !golfer.Earned("independence-day") {
		var (
			now   = time.Now().UTC()
			start = time.Date(2022, time.July, 4, 0, 0, 0, 0, time.UTC)
			end   = time.Date(2022, time.July, 5, 0, 0, 0, 0, time.UTC)
		)

		if now.Before(end) {
			banner := banner{Type: "info"}
			cheevo := config.CheevoByID["independence-day"]

			location := golfer.TimeZone
			if location == nil {
				location = time.UTC
			}

			banner.Body = template.HTML("The " + cheevo.Emoji + " <b>" +
				cheevo.Name + "</b> achievement will ")

			if start.Before(now) {
				banner.Body += "stop being available in " +
					pretty.Time(end.In(location)) + "."
			} else {
				banner.Body += "be available in " +
					pretty.Time(start.In(location)) + "."
			}

			banners = append(banners, banner)
		}
	}

	// Latest hole (if unsolved).
	if hole := recentHoles[0]; !golfer.Solved(hole.ID) {
		banners = append(banners, banner{
			Type: "info",
			Body: template.HTML(`The <a href="/` + hole.ID + `">` +
				hole.Name + "</a> hole is now live! Why not try and solve it?"),
		})
	}

	return
}
