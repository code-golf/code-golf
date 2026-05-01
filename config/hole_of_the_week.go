package config

import (
	"bytes"
	"cmp"
	"html/template"
	"math/rand/v2"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/code-golf/code-golf/db"
	"github.com/code-golf/code-golf/views"
	"github.com/lib/pq"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var firstWeek = time.Date(2026, time.January, 19, 0, 0, 0, 0, time.UTC).Unix()

var whitespace = regexp.MustCompile(`\s+`)

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

	var b bytes.Buffer
	if err := views.Render(&b, "hole-of-the-week", struct {
		holeOfTheWeek
		Week int
	}{
		holeOfTheWeek: hl,
		Week:          1 + int((thisWeek.Unix()-firstWeek)/(7*86_400)),
	}); err != nil {
		panic(err)
	}

	return template.HTML(b.String()), thisWeek
}

func HoleOfTheWeekText() string {
	bar, _ := HoleOfTheWeek()

	doc, err := html.Parse(strings.NewReader(string(bar)))
	if err != nil {
		panic(err)
	}

	var text strings.Builder
	for node := range doc.Descendants() {
		if node.Type == html.TextNode {
			if node.Parent.DataAtom == atom.Sup {
				// Convert superscript runes to the unicode equivalent.
				for _, r := range node.Data {
					switch r {
					case 'd':
						r = 'ᵈ'
					case 'h':
						r = 'ʰ'
					case 'n':
						r = 'ⁿ'
					case 's':
						r = 'ˢ'
					case 't':
						r = 'ᵗ'
					}
					text.WriteRune(r)
				}
			} else {
				text.WriteString(whitespace.ReplaceAllString(node.Data, " "))
			}
		}
	}

	return strings.TrimSpace(text.String())
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
