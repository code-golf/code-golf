package routes

import (
	"database/sql"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var validLangs = map[string]bool{
	"bash": true,
	"javascript": true,
	"perl": true,
	"perl6": true,
	"php": true,
	"python": true,
	"ruby": true,
}

func scores(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var hole string
	var langs []string
	var showDuplicates bool

	// Unless we have "/scores", tear apart the URL, and redirect if invalid.
	if strings.HasPrefix(r.URL.Path, "/scores/") {
		langs = strings.Split(r.URL.Path[8:], "/")

		url := "/scores"

		// Shift the first "lang" off if it's actually a hole.
		if len(langs) > 0 {
			if _, ok := preambles[ langs[0] ]; ok {
				hole, langs = langs[0], langs[1:]

				url += "/" + hole
			}
		}

		// Pop the last "lang" off if it's actually the showDuplicates flag.
		if len(langs) > 0 && langs[len(langs)-1] == "show-duplicates" {
			showDuplicates = true
			langs = langs[:len(langs)-1]
		}

		sort.Slice(langs, func(i, j int) bool { return langs[i] < langs[j] })

		// Avoid duplicate langs.
		seen := map[string]bool{}

		for _, lang := range langs {
			if validLangs[lang] && !seen[lang] {
				url += "/" + lang
				seen[lang] = true
			}
		}

		// No point in listing EVERY lang.
		if len(seen) == len(validLangs) {
			url = "/scores"

			if hole != "" {
				url += "/" + hole
			}
		}

		if showDuplicates && hole != "" {
			url += "/show-duplicates"
		}

		if r.URL.Path != url {
			http.Redirect(w, r, url, 301)
			return
		}
	}

	userID := printHeader(w, r, 200)

	w.Write([]byte(
		"<script async src=" + scoresJsPath +
			"></script><main id=scores><select id=hole><option value>All Holes",
	))

	for _, v := range [][]string{
		{"12-days-of-christmas", "12 Days of Christmas"},
		{"99-bottles-of-beer", "99 Bottles of Beer"},
		{"arabic-to-roman", "Arabic to Roman"},
		{"christmas-trees", "Christmas Trees"},
		{"emirp-numbers", "Emirp Numbers"},
		{"evil-numbers", "Evil Numbers"},
		{"fibonacci", "Fibonacci"},
		{"fizz-buzz", "Fizz Buzz"},
		{"happy-numbers", "Happy Numbers"},
		{"odious-numbers", "Odious Numbers"},
		{"pangram-grep", "Pangram Grep"},
		{"pascals-triangle", "Pascal's Triangle"},
		{"pernicious-numbers", "Pernicious Numbers"},
		{"prime-numbers", "Prime Numbers"},
		{"quine", "Quine"},
		{"roman-to-arabic", "Roman to Arabic"},
		{"sierpiński-triangle", "Sierpiński Triangle"},
		{"seven-segment", "Seven Segment"},
		{"spelling-numbers", "Spelling Numbers"},
		{"π", "π"},
		{"φ", "φ"},
		{"𝑒", "𝑒"},
		{"τ", "τ"},
	} {
		w.Write([]byte("<option "))
		if hole == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select><select id=lang><option value>All Langs"))

	for _, v := range [][]string{
		{"bash", "Bash"},
		{"javascript", "JavaScript"},
		{"perl", "Perl"},
		{"perl6", "Perl 6"},
		{"php", "PHP"},
		{"python", "Python"},
		{"ruby", "Ruby"},
	} {
		w.Write([]byte("<option "))
		if len(langs) > 0 && langs[0] == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select>"))

	if hole != "" {
		w.Write([]byte("<label><input type=checkbox"))
		if showDuplicates {
			w.Write([]byte(" checked"))
		}
		w.Write([]byte(">Show duplicate entries per person</label>"))
	}

	w.Write([]byte("<table class=scores>"))

	var concat, distinct, table, where string

	if hole != "" {
		where += " AND hole = '" + hole + "'"
		concat = "' class=', lang, '>', TO_CHAR(strokes, 'FM99,999'), '<td>(', TO_CHAR(score, 'FM9,999'), ' point', CASE WHEN score"
	} else {
		concat = "'>', TO_CHAR(score, 'FM9,999'), '<td>(', count, ' hole', CASE WHEN count"
	}

	if showDuplicates {
		table = "scored_leaderboard"
	} else {
		distinct = "DISTINCT ON (hole, user_id)"
		table    = "summed_leaderboard"
	}

	if len(langs) > 0 {
		where += " AND lang IN('" + strings.Join(langs, "','") + "')"
	}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT `+distinct+`
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id,
		         lang
		    FROM solutions
		   WHERE true`+where+`
		ORDER BY hole, user_id, LENGTH(code), submitted
		), scored_leaderboard AS (
		  SELECT hole,
		         ROUND(
		             (
		                 COUNT(*) OVER (PARTITION BY hole)
		                 -
		                 RANK() OVER (PARTITION BY hole ORDER BY strokes)
		                 +
		                 1
		             )
		             *
		             (
		                 100.0
		                 /
		                 COUNT(*) OVER (PARTITION BY hole)
		             )
		         ) score,
		         strokes,
		         submitted,
		         user_id,
		         lang
		    FROM leaderboard
		), summed_leaderboard AS (
		  SELECT user_id,
		         SUM(strokes)   strokes,
		         SUM(score)     score,
		         COUNT(*),
		         MAX(submitted) submitted,
		         STRING_AGG(CONCAT(lang), '') lang
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT CONCAT(
		             '<tr',
		             CASE WHEN user_id = $1 THEN ' class=me' END,
		             '><td>',
		             TO_CHAR(
		                 RANK() OVER (ORDER BY score DESC),
		                 'FM999"<sup>"th"</sup>"'
		             ),
		             '<td><img src="//avatars.githubusercontent.com/',
		             login,
		             '?s=26"><a href="/users/',
		             login,
		             '">',
		             login,
		             '</a><td',
		             `+concat+` > 1 THEN 's' END,
		             ')<td><time datetime=',
		             TO_CHAR(submitted, 'YYYY-MM-DD"T"HH24:MI:SS"Z>"FMDD Mon'),
		             '</time>'
		         )
		    FROM `+table+`
		    JOIN users on user_id = id
		ORDER BY score DESC, submitted`,
		userID,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var row sql.RawBytes

		if err := rows.Scan(&row); err != nil {
			panic(err)
		}

		w.Write(row)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
