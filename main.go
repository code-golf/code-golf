package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"syscall"
	"time"

	"github.com/code-golf/code-golf/github"
	"github.com/code-golf/code-golf/middleware"
	"github.com/code-golf/code-golf/routes"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	log.SetFlags(log.Ltime)

	rand.Seed(time.Now().UTC().UnixNano())

	db, err := sql.Open("postgres", "")
	if err != nil {
		panic(err)
	}

	dbBeta, err := sql.Open("postgres", "dbname=code-golf-beta")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RedirectHost,
		middleware.Public,
		middleware.RedirectSlashes,
		middleware.Compress(5),
		// middleware.Downtime,
		middleware.DatabaseHandler(db, dbBeta),
		middleware.GolferHandler,
	)

	r.NotFound(routes.NotFound)

	r.Get("/", routes.Index)
	r.Get("/{hole}", routes.Hole)
	r.Get("/about", routes.About)
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AdminArea)
		r.Get("/", routes.Admin)
		r.Get("/solutions", routes.AdminSolutions)
		r.Get("/solutions/run", routes.AdminSolutionsRun)
	})
	r.Get("/assets/{asset}", routes.Asset)
	r.Get("/callback", routes.Callback)
	r.Get("/callback/beta", routes.CallbackBeta)
	r.Get("/callback/dev", routes.CallbackDev)
	r.Get("/feeds/{feed}", routes.Feed)
	r.Route("/golfer", func(r chi.Router) {
		r.Use(middleware.GolferArea)
		r.Post("/cancel-delete", routes.GolferCancelDelete)
		r.Post("/delete", routes.GolferDelete)
		r.Get("/settings", routes.GolferSettings)
		r.Post("/settings", routes.GolferSettingsPost)
	})
	r.Route("/golfers/{name}", func(r chi.Router) {
		r.Use(middleware.GolferInfoHandler)
		r.Get("/", routes.GolferTrophies)
		r.Get("/achievements", routes.GolferAchievements)
		r.Get("/holes", routes.GolferHoles)
		r.Get("/holes/{scoring}", routes.GolferHoles)
	})
	r.Get("/ideas", routes.Ideas)
	r.Get("/log-out", routes.LogOut)
	r.Get("/random", routes.Random)
	r.Get("/recent", routes.Recent)
	r.Get("/recent/{lang}", routes.Recent)
	r.Get("/scores/{hole}/{lang}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/{scoring}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/{scoring}/{suffix}", routes.Scores)
	r.Get("/sitemap.xml", routes.Sitemap)
	r.Post("/solution", routes.Solution)
	r.Get("/stats", routes.Stats)
	r.Get("/users/{name}", routes.User)

	certManager := autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"code.golf", "beta.code.golf", "www.code.golf",
			// Legacy domain.
			"code-golf.io", "www.code-golf.io",
		),
	}

	server := &http.Server{
		Addr:    ":1443",
		Handler: r,
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // HTTP/2-required.
			},
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	var crt, key string
	if _, dev := syscall.Getenv("DEV"); dev {
		crt = "localhost.pem"
		key = "localhost-key.pem"
	} else {
		server.TLSConfig.GetCertificate = certManager.GetCertificate
	}

	// Every minute.
	go func() {
		for range time.NewTicker(time.Minute).C {
			for _, db := range []*sql.DB{db, dbBeta} {
				if _, err := db.Exec(
					"REFRESH MATERIALIZED VIEW CONCURRENTLY medals",
				); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	// Every 5 minutes.
	go func() {
		// Various GitHub API requests.
		for range time.NewTicker(5 * time.Minute).C {
			github.Run(db)
		}
	}()

	// Every hour.
	go func() {
		for range time.NewTicker(time.Hour).C {
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

			if err := syncToBeta(db, dbBeta); err != nil {
				log.Println(err)
			}
		}
	}()

	log.Println("Listening…")

	// Redirect HTTP to HTTPS and handle ACME challenges.
	go func() {
		panic(http.ListenAndServe(":1080", certManager.HTTPHandler(nil)))
	}()

	// Serve HTTPS.
	panic(server.ListenAndServeTLS(crt, key))
}

func syncToBeta(db, dbBeta *sql.DB) error {
	insertUser, err := dbBeta.Prepare(
		"INSERT INTO users(id, login) VALUES($1, $2) ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT id, login FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var login string

		if err := rows.Scan(&id, &login); err != nil {
			return err
		}

		if _, err := insertUser.Exec(id, login); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	insertCode, err := dbBeta.Prepare(
		"INSERT INTO code (code) VALUES ($1) ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}

	var insertBytes, insertChars *sql.Stmt

	for _, scoring := range []string{"bytes", "chars"} {
		stmt, err := dbBeta.Prepare(
			`WITH new_code AS (
				SELECT id, ` + scoring + ` FROM code WHERE code = $1
			) INSERT INTO solutions (code_id, failing, hole, lang, scoring, submitted, user_id)
				   SELECT id, $2, $3, $4, '` + scoring + `', $5, $6 FROM new_code
			  ON CONFLICT ON CONSTRAINT solutions_pkey
			DO UPDATE SET failing = excluded.failing,
						submitted = CASE
							WHEN (SELECT ` + scoring + ` FROM new_code)
							   < (SELECT ` + scoring + ` FROM code WHERE id = solutions.code_id)
							THEN excluded.submitted
							ELSE solutions.submitted
						END,
						  code_id = CASE
							WHEN (SELECT ` + scoring + ` FROM new_code)
							  <= (SELECT ` + scoring + ` FROM code WHERE id = solutions.code_id)
							THEN excluded.code_id
							ELSE solutions.code_id
						END`,
		)
		if err != nil {
			return err
		}

		if scoring == "bytes" {
			insertBytes = stmt
		} else {
			insertChars = stmt
		}
	}

	rows, err = db.Query(
		`SELECT code, failing, hole, lang, scoring, submitted, user_id
		   FROM solutions
		   JOIN code ON code_id = id`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var failing bool
		var code, hole, lang, scoring string
		var submitted time.Time
		var userID int

		if err := rows.Scan(
			&code,
			&failing,
			&hole,
			&lang,
			&scoring,
			&submitted,
			&userID,
		); err != nil {
			return err
		}

		if _, err := insertCode.Exec(code); err != nil {
			return err
		}

		if _, err := insertBytes.Exec(code, failing, hole, lang, submitted, userID); err != nil {
			return err
		}

		if _, err := insertChars.Exec(code, failing, hole, lang, submitted, userID); err != nil {
			return err
		}
	}

	return rows.Err()
}
