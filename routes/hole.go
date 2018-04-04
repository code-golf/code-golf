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
			"<script async src=" + holeJsPath + "></script><div id=status><div>" +
			"<h2>Program Arguments</h2><pre id=Arg></pre>" +
			"<h2>Standard Error</h2><pre id=Err></pre>" +
			"<h2>Expected Output</h2><pre id=Exp></pre>" +
			"<h2>Standard Output</h2><pre id=Out></pre>" +
			"</div></div><main id=hole",
	))

	if userID == 0 {
		w.Write([]byte(
			"><div id=alert>Please " +
				`<a href="//github.com/login/oauth/authorize?` +
				`client_id=7f6709819023e9215205&scope=user:email">` +
				"Login with GitHub</a> in order to save solutions " +
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

	w.Write([]byte(
		"><h1>" + hole.Name + "</h1><p>" + hole.Preamble +
			"<button>Run</button><div id=tabs>" +
			"<a href=#bash><div>Bash</div><div>not tried</div></a>" +
			"<a href=#javascript><div>JS</div><div>not tried</div></a>" +
			"<a href=#lisp><div>Lisp</div><div>not tried</div></a>" +
			"<a href=#lua><div>Lua</div><div>not tried</div></a>" +
			"<a href=#perl><div>Perl</div><div>not tried</div></a>" +
			"<a href=#perl6><div>Perl 6</div><div>not tried</div></a>" +
			"<a href=#php><div>PHP</div><div>not tried</div></a>" +
			"<a href=#python><div>Python</div><div>not tried</div></a>" +
			"<a href=#ruby><div>Ruby</div><div>not tried</div></a>" +
			"</div>",
	))
}
