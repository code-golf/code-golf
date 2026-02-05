package config

import (
	"cmp"
	"html/template"
	"math/rand/v2"
	"slices"
	"strings"
	"time"

	"github.com/code-golf/code-golf/db"
	"github.com/code-golf/code-golf/pretty"
	"github.com/lib/pq"
)

type holeOfTheWeek struct {
	Hole  *Hole
	Langs []*Lang
}

var holes = map[time.Time]holeOfTheWeek{}

func HoleOfTheWeek() (template.HTML, time.Time) {
	thisWeek := ThisWeek()
	hl, ok := holes[thisWeek]
	if !ok {
		return "", thisWeek
	}

	html := `<a href="/` + hl.Hole.ID + `">` + hl.Hole.Name +
		"</a> is the Hole of the Week. Solve it in either "

	for i, lang := range hl.Langs {
		if i > 0 {
			html += ", "

			if i == len(hl.Langs)-1 {
				html += "or "
			}
		}

		html += `<a href="/` + hl.Hole.ID + "#" + lang.ID + `">` + lang.Name + "</a>"
	}

	html += " within the next " + string(pretty.Time(thisWeek.AddDate(0, 0, 7))) +
		" to complete the challenge."

	return template.HTML(html), thisWeek
}

func PopulateHolesOfTheWeek(db db.Queryable) error {
	thisWeek := ThisWeek()
	nextWeek := thisWeek.AddDate(0, 0, 7)

	for _, week := range []time.Time{thisWeek, nextWeek} {
		hole := HoleList[rand.IntN(len(HoleList))]
		langs := make(Langs, 3)
		for i, j := range rand.Perm(len(LangList))[:len(langs)] {
			langs[i] = LangList[j]
		}

		slices.SortFunc(langs, func(a, b *Lang) int {
			return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		})

		if _, err := db.Exec(
			`INSERT INTO weekly_holes (week, hole, langs)
			      VALUES              (  $1,   $2,    $3)
			 ON CONFLICT DO NOTHING`,
			week, hole, pq.Array(langs),
		); err != nil {
			return err
		}

		langs = langs[:0]
		if err := db.QueryRow(
			"SELECT hole, langs FROM weekly_holes WHERE week = $1", week,
		).Scan(&hole, &langs); err != nil {
			return err
		}

		holes[week] = holeOfTheWeek{hole, langs}
	}

	return nil
}

// ThisWeek returns the start of the current week (Monday, NOT Sunday).
func ThisWeek() time.Time {
	year, month, day := time.Now().UTC().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today.AddDate(0, 0, -(int(today.Weekday())+6)%7)
}
