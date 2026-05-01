package routes

import (
	"net/http"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/oauth"
	"github.com/code-golf/code-golf/session"
)

// GET /golfers/{golfer}
func golferGET(w http.ResponseWriter, r *http.Request) {
	const limit = 100

	type row struct {
		AvatarURL, Date, Golfer string
		Cheevos                 []*config.Cheevo
		Hole                    *config.Hole
		Lang                    *config.Lang
	}

	db := session.Database(r)
	golferInfo := session.GolferInfo(r)

	location := time.UTC
	if golfer := session.Golfer(r); golfer != nil {
		location = golfer.Location()
	}

	data := struct {
		CategoryOverview map[string]map[string]int
		Connections      []oauth.Connection
		Followers        []golfer.GolferLink
		Following        []struct {
			golfer.GolferLink

			Bytes, Chars, Rank int
		}
		Langs []struct {
			Lang                          *config.Lang
			RankBytes, RankChars, RankMin *int
		}
		OAuthProviders map[string]*oauth.Config
		Wall           []row
	}{
		CategoryOverview: map[string]map[string]int{"bytes": {}, "chars": {}},
		Connections:      oauth.GetConnections(db, golferInfo.ID, true),
		OAuthProviders:   oauth.Providers,
		Wall:             make([]row, 0, limit),
	}

	// Note we hide followers as well as following if the bool is false.
	// Maybe this means the setting is badly named?
	rows, err := db.Query(
		`WITH data AS (
		 -- Cheevos
		    SELECT earned     date,
		           cheevo     cheevo,
		           NULL::hole hole,
		           NULL::lang lang,
		           user_id
		      FROM cheevos
		     WHERE CASE WHEN $1 THEN user_id = ANY(following($2, $3))
		                        ELSE user_id = $2
		           END
		 UNION ALL
		 -- Follows
		    SELECT followed     date,
		           NULL::cheevo cheevo,
		           NULL::hole   hole,
		           NULL::lang   lang,
		           follower_id  user_id
		      FROM follows
		     WHERE followee_id = $2 AND $1
		 UNION ALL
		 -- Holes
		    SELECT MAX(submitted) date,
		           NULL::cheevo   cheevo,
		           hole           hole,
		           lang           lang,
		           user_id
		      FROM solutions
		     WHERE CASE WHEN $1 THEN user_id = ANY(following($2, $3))
		                        ELSE user_id = $2
		           END
		  GROUP BY user_id, hole, lang
		) SELECT cheevo, date, avatar_url, name, hole, lang
		    FROM data JOIN golfers_with_avatars ON id = user_id
		ORDER BY date DESC, name, cheevo LIMIT $4`,
		session.Settings(r)["golfer/profile"]["followed-golfers-in-feed"],
		golferInfo.ID,
		golferInfo.FollowLimit(),
		limit,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

rows:
	for rows.Next() {
		var cheevo *config.Cheevo
		var time time.Time

		var r row
		if err := rows.Scan(
			&cheevo, &time, &r.AvatarURL, &r.Golfer, &r.Hole, &r.Lang,
		); err != nil {
			panic(err)
		}

		r.Date = time.In(location).Format("Mon 2 Jan 2006")

		if cheevo != nil {
			// Try and find a place in the current day to append the cheevo.
			for i := len(data.Wall) - 1; i >= 0 && data.Wall[i].Date == r.Date; i-- {
				if data.Wall[i].Cheevos != nil && data.Wall[i].Golfer == r.Golfer {
					data.Wall[i].Cheevos = append(data.Wall[i].Cheevos, cheevo)
					continue rows
				}
			}

			r.Cheevos = append(r.Cheevos, cheevo)
		}

		data.Wall = append(data.Wall, r)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	var categories []struct {
		Category, Scoring string
		Points            int
	}
	if err := db.Select(
		&categories,
		`WITH max_points_per_hole AS (
		    SELECT DISTINCT ON (hole, scoring) hole, scoring, points
		      FROM rankings
		     WHERE user_id = $1 AND NOT experimental
		  ORDER BY hole, scoring, points DESC
		) SELECT category,
		         ROUND(AVG((points)))   points,
		         scoring                scoring
		    FROM max_points_per_hole
		    JOIN holes ON id = hole
		GROUP BY scoring, category`,
		golferInfo.ID,
	); err != nil {
		panic(err)
	}

	// Ensure we have no missing categories.
	for _, hole := range config.HoleList {
		data.CategoryOverview["bytes"][hole.Category] = 0
		data.CategoryOverview["chars"][hole.Category] = 0
	}
	for _, cat := range categories {
		data.CategoryOverview[cat.Scoring][cat.Category] = cat.Points
	}

	if err := db.Select(
		&data.Langs,
		`WITH ranks AS (
		    SELECT user_id, scoring, lang,
		           RANK() OVER (PARTITION BY scoring, lang
		                            ORDER BY SUM(points_for_lang) DESC)
		      FROM rankings
		     WHERE NOT experimental
		  GROUP BY user_id, scoring, lang
		) SELECT lang,
		         any_value(rank) FILTER (WHERE scoring = 'bytes') rank_bytes,
		         any_value(rank) FILTER (WHERE scoring = 'chars') rank_chars,
		         min(rank)                                        rank_min
		    FROM ranks
		    JOIN langs ON lang = id
		   WHERE user_id = $1
		GROUP BY lang
		ORDER BY rank_min, any_value(name)`,
		golferInfo.ID,
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Following,
		`WITH follows AS (
		    SELECT avatar_url, country_flag, name,
		           COALESCE((SELECT points FROM points
		              WHERE scoring = 'bytes' AND user_id = id), 0) bytes,
		           COALESCE((SELECT points FROM points
		              WHERE scoring = 'chars' AND user_id = id), 0) chars
		      FROM golfers_with_avatars
		     WHERE id = ANY(following($1, $2))
		) SELECT *, RANK() OVER (ORDER BY bytes DESC, chars DESC)
		    FROM follows
		ORDER BY rank, name`,
		golferInfo.ID,
		golferInfo.FollowLimit(),
	); err != nil {
		panic(err)
	}

	if err := db.Select(
		&data.Followers,
		` SELECT avatar_url, name
		    FROM follows
		    JOIN golfers_with_avatars ON id = follower_id
		   WHERE followee_id = $1
		ORDER BY name`,
		golferInfo.ID,
	); err != nil {
		panic(err)
	}

	render(w, r, "golfer/profile", data, golferInfo.Name)
}
