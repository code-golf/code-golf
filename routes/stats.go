package routes

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var noHoles = strconv.Itoa(len(holes))
var noLangs = strconv.Itoa(len(langs))

func stats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	printHeader(w, r, 200)

	var data []byte

	w.Write([]byte(
		"<link rel=stylesheet href=" + statsCssPath + "><script async src=" +
			statsJsPath + "></script><main><div><div><span data-x=" + noLangs +
			">0</span>Languages</div></div><div><div><span data-x=" + noHoles +
			">0</span>Holes</div></div><div><div><span data-x=",
	))

	if err := db.QueryRow("SELECT COUNT(DISTINCT user_id) FROM solutions WHERE NOT failing").Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	w.Write([]byte(">0</span>Golfers</div></div><div><div><span data-x="))

	if err := db.QueryRow("SELECT COUNT(*) FROM solutions WHERE NOT failing").Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	// FIXME Make "Holes by Difficulty" data dynamic.
	w.Write([]byte(">0</span>Solutions</div></div><div><div><canvas data-data=[11,12,4]></canvas></div></div><div><div><canvas data-data="))

	if err := db.QueryRow(
		`SELECT ARRAY_TO_JSON(ARRAY_AGG(count))
		   FROM (SELECT COUNT(*) FROM solutions WHERE NOT failing GROUP BY lang ORDER BY lang) a`,
	).Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	w.Write([]byte("></canvas></div></div><div><div><canvas data-data='"))

	if err := db.QueryRow(
		`SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(t)))
		   FROM (SELECT x, COUNT(*) y
		   FROM (SELECT COUNT(DISTINCT hole) x FROM solutions WHERE NOT failing GROUP BY user_id) z
		GROUP BY x ORDER BY x) t`,
	).Scan(&data); err != nil {
		panic(err)
	} else {
		w.Write(data)
	}

	w.Write([]byte("'></canvas></div></div>"))
}
