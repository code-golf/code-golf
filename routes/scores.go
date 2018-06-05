package routes

import (
	"database/sql"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func scores(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var holeID string
	var langIDs []string
	var showDuplicates bool

	// Unless we have "/scores", tear apart the URL, and redirect if invalid.
	if strings.HasPrefix(r.URL.Path, "/scores/") {
		langIDs = strings.Split(r.URL.Path[8:], "/")

		url := "/scores"

		// Shift the first "lang" off if it's actually a hole.
		if len(langIDs) > 0 {
			if _, ok := holeByID[langIDs[0]]; ok {
				holeID, langIDs = langIDs[0], langIDs[1:]

				url += "/" + holeID
			}
		}

		// Pop the last "lang" off if it's actually the showDuplicates flag.
		if len(langIDs) > 0 && langIDs[len(langIDs)-1] == "show-duplicates" {
			showDuplicates = true
			langIDs = langIDs[:len(langIDs)-1]
		}

		sort.Slice(langIDs, func(i, j int) bool { return langIDs[i] < langIDs[j] })

		// Avoid duplicate langIDs.
		seen := map[string]bool{}

		for _, lang := range langIDs {
			if _, ok := langByID[lang]; ok && !seen[lang] {
				url += "/" + lang
				seen[lang] = true
			}
		}

		// No point in listing EVERY lang.
		if len(seen) == len(langs) {
			url = "/scores"

			if holeID != "" {
				url += "/" + holeID
			}
		}

		if showDuplicates && holeID != "" {
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

	for _, hole := range holes {
		w.Write([]byte("<option "))
		if holeID == hole.ID {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + hole.ID + ">" + hole.Name))
	}

	w.Write([]byte("</select><select id=lang><option value>All langs"))

	for _, lang := range langs {
		w.Write([]byte("<option "))
		if len(langIDs) > 0 && langIDs[0] == lang.ID {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + lang.ID + ">" + lang.Name))
	}

	w.Write([]byte("</select>"))

	if holeID != "" {
		w.Write([]byte("<label><input type=checkbox"))
		if showDuplicates {
			w.Write([]byte(" checked"))
		}
		w.Write([]byte(">Show duplicate entries per person</label>"))
	}

	w.Write([]byte("<table class=scores>"))

	var concat, distinct, table, where string

	if holeID != "" {
		where += " AND hole = '" + holeID + "'"
		concat = "' class=', lang, '>', TO_CHAR(strokes, 'FM99,999'), '<td>(', TO_CHAR(score, 'FM9,999'), ' point', CASE WHEN score"
	} else {
		concat = "'>', TO_CHAR(score, 'FM9,999'), '<td>(', count, ' hole', CASE WHEN count"
	}

	if showDuplicates {
		table = "scored_leaderboard"
	} else {
		distinct = "DISTINCT ON (hole, user_id)"
		table = "summed_leaderboard"
	}

	if len(langIDs) > 0 {
		where += " AND lang IN('" + strings.Join(langIDs, "','") + "')"
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
		   WHERE NOT failing`+where+`
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
		                 RANK() OVER (ORDER BY score DESC, strokes),
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
		ORDER BY score DESC, strokes, submitted`,
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
