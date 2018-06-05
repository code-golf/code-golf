package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func holeNg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	holeID := r.URL.Path[4:]
	userID := printHeader(w, r, 200)

	w.Write([]byte(
		"<link rel=stylesheet href=" + holeNgCssPath + ">" +
			"<script async src=" + holeNgJsPath + "></script><main id=hole",
	))

	if userID == 0 {
		w.Write([]byte(
			"><div id=alert>Please " +
				`<a href="//github.com/login/oauth/authorize?` +
				`client_id=7f6709819023e9215205&scope=user:email">` +
				"login with GitHub</a> in order to save solutions " +
				"and appear on the leaderboards.</div",
		))
	}

	hole := holeByID[holeID]

	w.Write([]byte(
		"><h1>" + hole.Name + "</h1><p>" + hole.Preamble + "<div id=tabs>"))

	for _, lang := range langs {
		w.Write([]byte(
			"<a href=#" + lang.ID + ` title="` +
				lang.Name + `">not tried</a>`,
		))
	}

	w.Write([]byte("</div>"))

	var html []byte

	// Fetch all the code per lang.
	if err := db.QueryRow(
		`SELECT STRING_AGG(CONCAT(
		              '<div class=code data-lang=',
		              unnest,
		              ' data-submitted="',
		              submitted,
		              '">',
		              REPLACE(code, '<', '&lt;'),
		              '</div>'
		          ), '' ORDER BY unnest)
		     FROM (SELECT UNNEST(ENUM_RANGE(NULL::lang))) z
		LEFT JOIN solutions
		       ON lang = unnest AND hole = $1 AND user_id = $2`,
		holeID, userID,
	).Scan(&html); err != nil {
		panic(err)
	} else {
		w.Write(html)
	}

	w.Write([]byte(
		"<button>Run</button><div id=status>" +
			"<h2></h2>" +
			"<aside id=arg><h3>Arg</h3><div></div></aside>" +
			"<aside id=err><h3>Err</h3><div></div></aside>" +
			"<aside id=exp><h3>Exp</h3><div></div></aside>" +
			"<aside id=out><h3>Out</h3><div></div></aside>" +
			"</div>",
	))
}
