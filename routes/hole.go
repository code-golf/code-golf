package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func hole(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	holeID := r.URL.Path[1:]
	userID := printHeader(w, r, 200)

	w.Write([]byte(
		"<link rel=stylesheet href=" + holeCssPath + ">" +
			"<script async src=" + holeJsPath + "></script><main id=hole",
	))

	if userID == 0 {
		w.Write([]byte(
			"><div id=alert>Please " +
				`<a href="//github.com/login/oauth/authorize?` +
				`client_id=7f6709819023e9215205&scope=user:email">` +
				"login with GitHub</a> in order to save solutions " +
				"and appear on the leaderboards.</div",
		))
	} else {
		var html []byte

		// Fetch the latest successful lang.
		if err := db.QueryRow(
			`SELECT CONCAT(' data-lang=', lang)
			   FROM solutions
			  WHERE user_id = $1 AND hole = $2`,
			userID, holeID,
		).Scan(&html); err == nil {
			w.Write(html)
		} else if err != sql.ErrNoRows {
			panic(err)
		}

		// Fetch all the code per lang.
		if err := db.QueryRow(
			`SELECT STRING_AGG(CONCAT(
			            ' data-',
			            lang,
			            '="',
			            REPLACE(code, '"', '&#34;'),
			            '"'
			        ), '')
			   FROM solutions
			  WHERE user_id = $1 AND hole = $2`,
			userID, holeID,
		).Scan(&html); err == nil {
			w.Write(html)
		} else if err != sql.ErrNoRows {
			panic(err)
		}
	}

	hole := holeByID[holeID]

	w.Write([]byte("><nav><a href=" + hole.Prev + ">Previous Hole</a><a href=random>Random Hole</a><a href=" + hole.Next + ">Next Hole</a></nav"))

	w.Write([]byte(
		"><h1>" + hole.Name + "</h1><div><p>" + hole.Preamble + "</div><div id=tabs>"))

	for _, l := range langs {
		w.Write([]byte("<a href=#" + l.ID + ` title="` + l.Name + `"></a>`))
	}

	w.Write([]byte(
		"</div><div id=wrapper></div>" +
			`<div class="info javascript"><b>arguments</b> holds ARGV, <b>print()</b> to output with a newline, <b>write()</b> to output without a newline.</div>` +
			`<div class="info perl"><b>say</b>, <b>signatures</b>, and <b>state</b> are available without any import.</div>` +
			"<button>Run</button>" +
			"<div id=status>" +
			"<h2></h2>" +
			"<aside id=err><h3>Errors</h3><div></div></aside>" +
			"<aside id=arg><h3>Arguments</h3><div></div></aside>" +
			"<aside id=exp><h3>Expected</h3><div></div></aside>" +
			"<aside id=out><h3>Output</h3><div></div></aside>" +
			"</div>",
	))
}
