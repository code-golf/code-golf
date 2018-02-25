package routes

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var noHoles = strconv.Itoa(len(holes))
var noLangs = strconv.Itoa(len(validLangs))

func stats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	printHeader(w, r, 200)

	var data []byte

	w.Write([]byte(
		"<link rel=stylesheet href=" + statsCssPath + "><script async src=" +
			statsJsPath + "></script><main><div><div><span>" + noLangs +
			"</span>Languages</div></div><div><div><span>" + noHoles +
			"</span>Holes</div></div><div><div><span>",
	))

	if err := db.QueryRow("SELECT COUNT(DISTINCT user_id) FROM solutions").Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	w.Write([]byte("</span>Golfers</div></div><div><div><span>"))

	if err := db.QueryRow("SELECT TO_CHAR(COUNT(*), 'FM9,999') FROM solutions").Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	// FIXME Make "Holes by Difficulty" data dynamic.
	w.Write([]byte("</span>Solutions</div></div><div><div><canvas data-data=[11,12,4]></canvas></div></div><div><div><canvas data-data="))

	if err := db.QueryRow(
		`SELECT ARRAY_TO_JSON(ARRAY_AGG(count))
		   FROM (SELECT COUNT(*) FROM solutions GROUP BY lang ORDER BY lang) a`,
	).Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	w.Write([]byte("></canvas></div></div><div><div><canvas data-data='"))

	if err := db.QueryRow(
		`SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(t)))
		   FROM (SELECT CASE WHEN x = 1 THEN 1.05 ELSE x END, COUNT(*) y
		   FROM (SELECT COUNT(DISTINCT hole) x FROM solutions GROUP BY user_id) z
		GROUP BY x ORDER BY x) t`,
	).Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	w.Write([]byte("'></canvas></div></div>"))
}
