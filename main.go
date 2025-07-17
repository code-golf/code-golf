package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/db"
	"github.com/code-golf/code-golf/discord"
	"github.com/code-golf/code-golf/github"
	"github.com/code-golf/code-golf/routes"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.Ltime)

	db := db.Open()

	// Attempt to populate the holes/langs tables every sec until we succeed.
	// This handles the site starting before the DB.
	go func() {
		if err := populateHolesLangsTables(db); err != nil {
			log.Println(err)
		} else {
			return
		}

		for range time.Tick(time.Second) {
			if err := populateHolesLangsTables(db); err != nil {
				log.Println(err)
			} else {
				break
			}
		}
	}()

	// Every 10 seconds.
	go func() {
		// Refreshing the mat views every 10 seconds is overkill on dev.
		// FIXME Maybe it would be better to do something with NOTIFY/TRIGGER.
		duration := 10 * time.Second
		if _, dev := os.LookupEnv("DEV"); dev {
			duration = 5 * time.Minute
		}

		for range time.Tick(duration) {
			for _, view := range []string{"medals", "rankings", "points"} {
				if _, err := db.Exec(
					"REFRESH MATERIALIZED VIEW CONCURRENTLY " + view,
				); err != nil {
					log.Println(err)
				}
			}

			// Once points is refreshed, award points based cheevos.
			if _, err := db.Exec(
				`INSERT INTO trophies(user_id, trophy)
				      SELECT user_id, 'big-brother'::cheevo
				        FROM points
				       WHERE points >= 1_984
				   UNION ALL
				      SELECT user_id, 'its-over-9000'
				        FROM points
				       WHERE points > 9_000
				   UNION ALL
				      SELECT user_id, 'twenty-kiloleagues'
				        FROM points
				       WHERE points >= 20_000
				   UNION ALL
				      SELECT user_id, 'marathon-runner'
				        FROM points
				       WHERE points >= 42_195
				   UNION ALL
				      SELECT user_id, '0xdead'
				        FROM points
				       WHERE points >= 57_005
				   UNION ALL
				      SELECT user_id, 'overflowing'
				        FROM points
				       WHERE points >= 65_536
				   UNION ALL
				      SELECT user_id, 'into-space'
				        FROM points
				       WHERE points >= 100_000
				 ON CONFLICT DO NOTHING`,
			); err != nil {
				log.Println(err)
			}
		}
	}()

	// Every 5 minutes.
	go func() {
		for range time.Tick(5 * time.Minute) {
			// Various GitHub API requests.
			github.Run(db, false)

			if err := github.Wiki(db); err != nil {
				log.Println(err)
			}

			if err := discord.AwardRoles(db); err != nil {
				log.Println(err)
			}
		}
	}()

	// Every hour.
	go func() {
		for range time.Tick(time.Hour) {
			// Update GitHub usernames.
			github.Run(db, true)

			for _, job := range [...]struct{ name, sql string }{
				{
					"expired sessions",
					`DELETE FROM sessions
					  WHERE last_used < TIMEZONE('UTC', NOW()) - INTERVAL '30 days'`,
				},
				{
					"superfluous users",
					`DELETE FROM users u
					  WHERE NOT EXISTS (SELECT FROM sessions WHERE user_id = u.id)
						AND NOT EXISTS (SELECT FROM trophies WHERE user_id = u.id)`,
				},
				{
					"users scheduled for deletion",
					"DELETE FROM users WHERE delete < TIMEZONE('UTC', NOW())",
				},
			} {
				if res, err := db.Exec(job.sql); err != nil {
					log.Println(err)
				} else if rows, _ := res.RowsAffected(); rows != 0 {
					log.Printf("Deleted %d %s\n", rows, job.name)
				}
			}

			if _, err := db.Exec(
				`INSERT INTO trophies(user_id, trophy)
				 SELECT user_id, 'aged-like-fine-wine'
				   FROM solutions
				  WHERE NOT failing
				  GROUP BY user_id
				 HAVING EXTRACT(days FROM TIMEZONE('UTC', NOW()) - MIN(submitted)) >= 365
					 ON CONFLICT DO NOTHING`,
			); err != nil {
				log.Println(err)
			}
		}
	}()

	log.Println("Listening…")

	// Dev.
	if _, dev := os.LookupEnv("DEV"); dev {
		// Redirect HTTP to HTTPS.
		go func() {
			panic(http.ListenAndServe(":80",
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, "https://localhost"+r.RequestURI,
						http.StatusMovedPermanently)
				})))
		}()

		// Serve HTTPS.
		panic(http.ListenAndServeTLS(
			":443", "localhost.pem", "localhost-key.pem", routes.Router(db)))
	}

	// Live only listens on HTTP, TLS is handled by Caddy.
	panic(http.ListenAndServe(":80", routes.Router(db)))
}

func populateHolesLangsTables(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertHole, err := tx.PrepareNamed(
		`INSERT INTO holes ( id,  experiment)
		      VALUES       (:id, :experiment)`,
	)
	if err != nil {
		return err
	}

	if _, err := tx.Exec("TRUNCATE holes"); err != nil {
		return err
	}

	for _, hole := range config.AllHoleList {
		if _, err := insertHole.Exec(hole); err != nil {
			return err
		}
	}

	insertLang, err := tx.PrepareNamed(
		`INSERT INTO langs ( id,  experiment,  digest_trunc,  name)
		      VALUES       (:id, :experiment, :digest_trunc, :name)`,
	)
	if err != nil {
		return err
	}

	if _, err := tx.Exec("TRUNCATE langs"); err != nil {
		return err
	}

	for _, lang := range config.AllLangList {
		if _, err := insertLang.Exec(lang); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Populated holes & langs tables.")
	return nil
}
