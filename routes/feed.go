package routes

import (
	"html"
	"net/http"
	"time"

	"github.com/code-golf/code-golf/hole"
	"github.com/code-golf/code-golf/lang"
	"github.com/gorilla/feeds"
)

var (
	atomFeed, jsonFeed, rssFeed []byte
	feed                        feeds.Feed
)

// TZ=UTC git log --date='format-local:%Y-%m-%d %X' --format='%h %cd %s'
func init() {
	feed = feeds.Feed{
		Link:  &feeds.Link{Href: "https://code.golf/"},
		Title: "Code Golf",
	}

	for _, i := range []struct {
		sha, created, id string
		hole             bool
	}{
		{"e3ade9f", "2021-05-03 20:17:00", "star-wars-opening-crawl", true},
		{"7b70234", "2021-05-01 04:08:12", "van-eck-sequence", true},
		{"db4bfba", "2021-04-25 22:39:11", "sudoku-v2", true},
		{"a3cbf07", "2021-03-23 02:10:24", "crystal", false},
		{"78f7023", "2021-01-24 20:44:04", "hexagony", false},
		{"1d9ce0d", "2021-01-16 22:57:45", "kolakoski-constant", true},
		{"1d9ce0d", "2021-01-16 22:57:45", "kolakoski-sequence", true},
		{"1ef979e", "2020-12-26 21:13:30", "recamán", true},
		{"74bbb70", "2020-12-24 23:44:45", "v", false},
		{"fe9bf4c", "2020-12-01 18:53:11", "look-and-say", true},
		{"c38e23e", "2020-11-25 19:52:33", "emojify", true},
		{"89fe682", "2020-11-11 01:20:43", "levenshtein-distance", true},
		{"4405e02", "2020-10-31 23:52:10", "fish", false},
		{"96034b1", "2020-10-30 16:57:20", "vampire-numbers", true},
		{"e3e56fe", "2020-10-24 21:58:36", "zig", false},
		{"3858950", "2020-10-04 23:09:53", "intersection", true},
		{"5a831a7", "2020-09-29 01:09:27", "tongue-twisters", true},
		{"9f96619", "2020-09-29 00:56:09", "sql", false},
		{"90b0cc8", "2020-07-26 01:05:14", "leyland-numbers", true},
		{"681846d", "2020-07-24 21:39:33", "cobol", false},
		{"3e971cc", "2020-06-21 01:44:16", "css-colors", true},
		{"5a2068c", "2020-06-01 04:28:26", "powershell", false},
		{"b71f2ee", "2020-05-09 13:17:53", "c-sharp", false},
		{"b71f2ee", "2020-05-09 13:17:53", "f-sharp", false},
		{"6b791f5", "2020-04-21 15:18:00", "java", false},
		{"4849dde", "2020-04-17 18:56:56", "fortran", false},
		{"5feb18b", "2020-04-15 17:33:16", "go", false},
		{"d386360", "2020-04-14 01:20:14", "lucky-tickets", true},
		{"08e4756", "2020-01-28 13:38:00", "united-states", true},
		{"93d765b", "2020-01-12 13:26:11", "swift", false},
		{"a9bbba9", "2020-01-03 19:05:17", "rust", false},
		{"fb227ed", "2019-11-17 17:32:14", "abundant-numbers", true},
		{"4e38900", "2019-07-11 16:34:13", "ordinal-numbers", true},
		{"de9369a", "2019-06-17 16:33:28", "rock-paper-scissors-spock-lizard", true},
		{"5348743", "2019-06-12 22:10:04", "ten-pin-bowling", true},
		{"14290db", "2019-05-25 10:22:43", "√2", true},
		{"269ff68", "2019-05-19 20:13:21", "nim", false},
		{"3cbc7cf", "2019-04-07 14:33:47", "brainfuck", false},
		{"c960e70", "2019-02-28 21:27:25", "c", false},
		{"c689b8a", "2019-01-13 17:09:36", "cubes", true},
		{"2f87dea", "2018-12-03 18:34:16", "leap-years", true},
		{"1178818", "2018-11-15 19:16:17", "sudoku", true},
		{"447121b", "2018-09-02 16:27:32", "julia", false},
		{"00660df", "2018-08-08 23:11:21", "j", false},
		{"edd2828", "2018-07-25 19:50:07", "poker", true},
		{"646df41", "2018-07-06 21:49:47", "rule-110", true},
		{"2080b94", "2018-06-07 07:22:01", "λ", true},
		{"834750b", "2018-06-06 20:54:22", "diamonds", true},
		{"bd8e789", "2018-05-20 18:09:59", "haskell", false},
		{"7b72ebc", "2018-05-03 17:30:42", "niven-numbers", true},
		{"827599e", "2018-03-22 16:56:44", "lisp", false},
		{"5790715", "2018-02-18 21:01:24", "morse-decoder", true},
		{"5790715", "2018-02-18 21:01:24", "morse-encoder", true},
		{"05e21ff", "2018-01-27 23:39:06", "brainfuck", true},
		{"922fb91", "2018-01-10 20:36:47", "divisors", true},
		{"d83fcf9", "2018-01-07 12:52:41", "lua", false},
		{"079513e", "2017-12-08 17:12:30", "12-days-of-christmas", true},
		{"30fc7c2", "2017-12-05 15:27:34", "christmas-trees", true},
		{"2dfbcfe", "2017-11-30 19:23:00", "pangram-grep", true},
		{"31d18c8", "2017-11-28 20:31:07", "τ", true},
		{"a3ef71c", "2017-11-12 03:11:36", "bash", false},
		{"0219147", "2017-11-11 00:02:40", "φ", true},
		{"63510fc", "2017-11-10 23:21:35", "roman-to-arabic", true},
		{"ee1742f", "2017-11-07 21:56:03", "quine", true},
		{"ce2d6b9", "2017-10-31 04:39:43", "happy-numbers", true},
		{"63bad53", "2017-10-18 19:13:56", "pernicious-numbers", true},
		{"e71743f", "2017-10-18 18:24:51", "evil-numbers", true},
		{"e71743f", "2017-10-18 18:24:51", "odious-numbers", true},
		{"35da66a", "2017-10-08 02:11:33", "sierpiński-triangle", true},
		{"aa9e81e", "2017-10-08 01:39:39", "emirp-numbers", true},
		{"ac9b179", "2017-10-04 13:00:28", "prime-numbers", true},
		{"39ce198", "2017-09-29 20:41:30", "spelling-numbers", true},
		{"bb1a117", "2017-09-20 17:18:23", "𝑒", true},
		{"7475f08", "2017-09-16 13:57:43", "fibonacci", true},
		{"7d34727", "2017-08-27 23:41:32", "seven-segment", true},
		{"b1d91b8", "2017-07-22 22:22:25", "arabic-to-roman", true},
		{"07aa3ed", "2017-07-16 00:18:21", "javascript", false},
		{"b1d083d", "2017-07-15 21:24:09", "python", false},
		{"26cf869", "2017-07-15 20:28:52", "pascals-triangle", true},
		{"15bc065", "2017-07-06 21:37:43", "99-bottles-of-beer", true},
		{"9cc775e", "2017-07-02 17:00:43", "π", true},
		{"dc1b9a8", "2017-06-14 22:03:21", "php", false},
		{"c5468f0", "2017-06-12 23:34:47", "raku", false},
		{"8029a96", "2017-05-08 23:06:22", "ruby", false},
	} {
		var name string
		if i.hole {
			name = hole.ByID[i.id].Name
		} else {
			name = lang.ByID[i.id].Name
		}

		var link string
		if i.hole {
			link = "https://code.golf/" + i.id
		} else {
			link = "https://code.golf/rankings/holes/all/" + i.id + "/bytes"
		}

		item := feeds.Item{
			Description: "Added the <a href=" + link + ">“" + html.EscapeString(name) + "”</a> ",
			Id:          link,
			Link:        &feeds.Link{Href: link},
			Title:       "Added “" + name + "” ",
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

// Feeds serves /feeds
func Feeds(w http.ResponseWriter, r *http.Request) {
	render(w, r, "feeds", feed, "Feeds")
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
