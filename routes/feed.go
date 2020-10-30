package routes

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

var atomFeed, jsonFeed, rssFeed []byte

// TZ=UTC git log --date='format-local:%Y-%m-%d %X' --format='%h %cd %s'
func init() {
	feed := feeds.Feed{
		Link:  &feeds.Link{Href: "https://code.golf/"},
		Title: "Code Golf",
	}

	for _, i := range []struct {
		sha, created, name string
		hole               bool
	}{
		{"e3e56fe", "2020-10-24 21:58:36", "Zig", false},
		{"3858950", "2020-10-04 23:09:53", "Intersection", true},
		{"5a831a7", "2020-09-29 01:09:27", "Tongue-twisters", true},
		{"9f96619", "2020-09-29 00:56:09", "SQL", false},
		{"90b0cc8", "2020-07-26 01:05:14", "Leyland Numbers", true},
		{"681846d", "2020-07-24 21:39:33", "COBOL", false},
		{"3e971cc", "2020-06-21 01:44:16", "CSS Colors", true},
		{"5a2068c", "2020-06-01 04:28:26", "PowerShell", false},
		{"b71f2ee", "2020-05-09 13:17:53", "C#", false},
		{"b71f2ee", "2020-05-09 13:17:53", "F#", false},
		{"6b791f5", "2020-04-21 15:18:00", "Java", false},
		{"4849dde", "2020-04-17 18:56:56", "Fortran", false},
		{"5feb18b", "2020-04-15 17:33:16", "Go", false},
		{"d386360", "2020-04-14 01:20:14", "Lucky Tickets", true},
		{"08e4756", "2020-01-28 13:38:00", "United States", true},
		{"93d765b", "2020-01-12 13:26:11", "Swift", false},
		{"a9bbba9", "2020-01-03 19:05:17", "Rust", false},
		{"fb227ed", "2019-11-17 17:32:14", "Abundant Numbers", true},
		{"4e38900", "2019-07-11 16:34:13", "Ordinal Numbers", true},
		{"de9369a", "2019-06-17 16:33:28", "Rock-paper-scissors-Spock-lizard", true},
		{"5348743", "2019-06-12 22:10:04", "Ten-pin Bowling", true},
		{"14290db", "2019-05-25 10:22:43", "√2", true},
		{"269ff68", "2019-05-19 20:13:21", "Nim", false},
		{"3cbc7cf", "2019-04-07 14:33:47", "brainfuck", false},
		{"c960e70", "2019-02-28 21:27:25", "C", false},
		{"c689b8a", "2019-01-13 17:09:36", "Cubes", true},
		{"2f87dea", "2018-12-03 18:34:16", "Leap Years", true},
		{"1178818", "2018-11-15 19:16:17", "Sudoku", true},
		{"447121b", "2018-09-02 16:27:32", "Julia", false},
		{"00660df", "2018-08-08 23:11:21", "J", false},
		{"edd2828", "2018-07-25 19:50:07", "Poker", true},
		{"646df41", "2018-07-06 21:49:47", "Rule 110", true},
		{"2080b94", "2018-06-07 07:22:01", "λ", true},
		{"834750b", "2018-06-06 20:54:22", "Diamonds", true},
		{"bd8e789", "2018-05-20 18:09:59", "Haskell", false},
		{"7b72ebc", "2018-05-03 17:30:42", "Niven Numbers", true},
		{"827599e", "2018-03-22 16:56:44", "Lisp", false},
		{"5790715", "2018-02-18 21:01:24", "Morse Decoder", true},
		{"5790715", "2018-02-18 21:01:24", "Morse Encoder", true},
		{"05e21ff", "2018-01-27 23:39:06", "brainfuck", true},
		{"922fb91", "2018-01-10 20:36:47", "Divisors", true},
		{"d83fcf9", "2018-01-07 12:52:41", "Lua", false},
		{"079513e", "2017-12-08 17:12:30", "12 Days of Christmas", true},
		{"30fc7c2", "2017-12-05 15:27:34", "Christmas Trees", true},
		{"2dfbcfe", "2017-11-30 19:23:00", "Pangram Grep", true},
		{"31d18c8", "2017-11-28 20:31:07", "τ", true},
		{"a3ef71c", "2017-11-12 03:11:36", "Bash", false},
		{"0219147", "2017-11-11 00:02:40", "φ", true},
		{"63510fc", "2017-11-10 23:21:35", "Roman to Arabic", true},
		{"ee1742f", "2017-11-07 21:56:03", "Quine", true},
		{"ce2d6b9", "2017-10-31 04:39:43", "Happy Numbers", true},
		{"63bad53", "2017-10-18 19:13:56", "Pernicious Numbers", true},
		{"e71743f", "2017-10-18 18:24:51", "Evil Numbers", true},
		{"e71743f", "2017-10-18 18:24:51", "Odious Numbers", true},
		{"35da66a", "2017-10-08 02:11:33", "Sierpiński Triangle", true},
		{"aa9e81e", "2017-10-08 01:39:39", "Emirp Numbers", true},
		{"ac9b179", "2017-10-04 13:00:28", "Prime Numbers", true},
		{"39ce198", "2017-09-29 20:41:30", "Spelling Numbers", true},
		{"bb1a117", "2017-09-20 17:18:23", "𝑒", true},
		{"7475f08", "2017-09-16 13:57:43", "Fibonacci", true},
		{"7d34727", "2017-08-27 23:41:32", "Seven Segment", true},
		{"b1d91b8", "2017-07-22 22:22:25", "Arabic to Roman", true},
		{"07aa3ed", "2017-07-16 00:18:21", "JavaScript", false},
		{"b1d083d", "2017-07-15 21:24:09", "Python", false},
		{"26cf869", "2017-07-15 20:28:52", "Pascal’s Triangle", true},
		{"15bc065", "2017-07-06 21:37:43", "99 Bottles of Beer", true},
		{"9cc775e", "2017-07-02 17:00:43", "π", true},
		{"dc1b9a8", "2017-06-14 22:03:21", "PHP", false},
		{"c5468f0", "2017-06-12 23:34:47", "Raku", false},
		{"8029a96", "2017-05-08 23:06:22", "Ruby", false},
	} {
		link := "https://code.golf/"

		if !i.hole {
			link += "scores/all-holes/"
		}

		link += url.PathEscape(
			strings.ReplaceAll(
				strings.ReplaceAll(strings.ToLower(i.name), " ", "-"),
				"#", "-sharp",
			),
		)

		item := feeds.Item{
			Description: "Added the <a href=" + link + ">“" + i.name + "”</a> ",
			Id:          link,
			Link:        &feeds.Link{Href: link},
			Title:       "Added “" + i.name + "” ",
		}

		if i.hole {
			item.Title += "Hole"
			item.Description += "hole"
		} else {
			item.Title += "Language"
			item.Description += "language"
		}

		item.Description += " via <a href=https://github.com/code-golf/code-golf/commit/" +
			i.sha + ">" + i.sha + "</a>."

		var err error
		if item.Created, err = time.Parse("2006-01-02 15:04:05", i.created); err != nil {
			panic(err)
		}

		feed.Items = append(feed.Items, &item)

		if feed.Created.IsZero() {
			feed.Created = item.Created
		}
	}

	feed.Title = "Code Golf (Atom Feed)"

	if data, err := feed.ToAtom(); err != nil {
		panic(err)
	} else {
		atomFeed = []byte(data)
	}

	feed.Title = "Code Golf (JSON Feed)"

	if data, err := feed.ToJSON(); err != nil {
		panic(err)
	} else {
		jsonFeed = []byte(data)
	}

	feed.Title = "Code Golf (RSS Feed)"

	if data, err := feed.ToRss(); err != nil {
		panic(err)
	} else {
		rssFeed = []byte(data)
	}
}

// Feed serves /feeds/{feed}
func Feed(w http.ResponseWriter, r *http.Request) {
	switch param(r, "feed") {
	case "atom":
		w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
		w.Write(atomFeed)
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonFeed)
	case "rss":
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		w.Write(rssFeed)
	default:
		NotFound(w, r)
	}
}
