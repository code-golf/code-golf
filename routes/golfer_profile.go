package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/oauth"
	"github.com/code-golf/code-golf/session"
	"github.com/lib/pq"
)

// GET /golfers/{golfer}
func golferGET(w http.ResponseWriter, r *http.Request) {
	const limit = 100

	type row struct {
		Cheevo *config.Cheevo
		Date   time.Time
		Golfer string
		Hole   *config.Hole
		Lang   *config.Lang
	}

	db := session.Database(r)
	golfer := session.GolferInfo(r)

	data := struct {
		Connections []oauth.Connection
		Followers   []string
		Following   []struct {
			Bytes, Chars, Rank int
			Country            config.NullCountry
			Name               string
		}
		Langs          []*config.Lang
		OAuthProviders map[string]*oauth.Config
		Trophies       map[string]map[string]int
		Wall           []row
	}{
		Connections:    oauth.GetConnections(db, golfer.ID, true),
		Langs:          config.LangList,
		OAuthProviders: oauth.Providers,
		Trophies:       map[string]map[string]int{},
		Wall:           make([]row, 0, limit),
	}

	rows, err := db.Query(
		`WITH data AS (
		 -- Cheevos
		    SELECT earned       date,
		           trophy::text cheevo,
		           ''           hole,
		           ''           lang,
		           user_id
		      FROM trophies
		     WHERE user_id = ANY(following($1, $2))
		 UNION ALL
		 -- Follows
		    SELECT followed    date,
		           ''          cheevo,
		           ''          hole,
		           ''          lang,
		           follower_id user_id
		      FROM follows
		     WHERE followee_id = $1
		 UNION ALL
		 -- Holes
		    SELECT MAX(submitted) date,
		           ''             cheevo,
		           hole::text     hole,
		           lang::text     lang,
		           user_id
		      FROM solutions
		     WHERE user_id = ANY(following($1, $2))
		  GROUP BY user_id, hole, lang
		) SELECT cheevo, date, login, hole, lang
		    FROM data JOIN users ON id = user_id
		ORDER BY date DESC, login LIMIT $3`,
		golfer.ID,
		golfer.FollowLimit(),
		limit,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cheevo, golfer, hole, lang string
		var date time.Time

		if err := rows.Scan(&cheevo, &date, &golfer, &hole, &lang); err != nil {
			panic(err)
		}

		data.Wall = append(data.Wall, row{
			Cheevo: config.CheevoByID[cheevo],
			Date:   date,
			Golfer: golfer,
			Hole:   config.HoleByID[hole],
			Lang:   config.LangByID[lang],
		})
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	rows, err = db.Query(
		`WITH ranks AS (
		    SELECT user_id, scoring, lang,
		           RANK() OVER (PARTITION BY scoring, lang
		                            ORDER BY SUM(points_for_lang) DESC)
		      FROM rankings
		  GROUP BY user_id, scoring, lang
		) SELECT lang, scoring, rank FROM ranks WHERE user_id = $1`,
		golfer.ID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var lang, scoring string
		var rank int

		if err := rows.Scan(&lang, &scoring, &rank); err != nil {
			panic(err)
		}

		if _, ok := data.Trophies[lang]; !ok {
			data.Trophies[lang] = map[string]int{}
		}

		data.Trophies[lang][scoring] = rank
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Following,
		`WITH follows AS (
		    SELECT country_flag country, login name,
		           COALESCE((SELECT points FROM points
		              WHERE scoring = 'bytes' AND user_id = id), 0) bytes,
		           COALESCE((SELECT points FROM points
		              WHERE scoring = 'chars' AND user_id = id), 0) chars
		      FROM users
		     WHERE id = ANY(following($1, $2))
		) SELECT *, RANK() OVER (ORDER BY bytes DESC, chars DESC)
		    FROM follows
		ORDER BY rank, name`,
		golfer.ID,
		golfer.FollowLimit(),
	); err != nil {
		panic(err)
	}

	if err := db.QueryRow(
		`SELECT array_agg(login ORDER BY login)
		   FROM follows
		   JOIN users ON id = follower_id
		  WHERE followee_id = $1`,
		golfer.ID,
	).Scan(pq.Array(&data.Followers)); err != nil {
		panic(err)
	}

	render(w, r, "golfer/profile", data, golfer.Name)
}
