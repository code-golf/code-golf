package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/code-golf/code-golf/github"
	"github.com/code-golf/code-golf/middleware"
	"github.com/code-golf/code-golf/routes"
	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RedirectHost,
		middleware.Static,
		middleware.RedirectSlashes,
		middleware.Compress(5),
		// middleware.Downtime,
		middleware.DatabaseHandler(db),
		middleware.GolferHandler,
	)

	r.NotFound(routes.NotFound)

	r.Get("/", routes.Index)
	r.Get("/{hole}", routes.Hole)
	r.Get("/ng/{hole}", routes.HoleNG)
	r.Get("/about", routes.About)
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AdminArea)
		r.Get("/", routes.Admin)
		r.Get("/solutions", routes.AdminSolutions)
		r.Get("/solutions/run", routes.AdminSolutionsRun)
	})
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.API)
		r.Get("/", routes.API)
		r.Get("/langs", routes.APILangs)
		r.Get("/langs/{lang}", routes.APILang)
		r.Get("/suggestions/golfers", routes.APISuggestionsGolfers)
	})
	r.Get("/callback", routes.Callback)
	r.Get("/callback/dev", routes.CallbackDev)
	r.Get("/feeds", routes.Feeds)
	r.Get("/feeds/{feed}", routes.Feed)
	r.Route("/golfer", func(r chi.Router) {
		r.Use(middleware.GolferArea)
		r.Post("/cancel-delete", routes.GolferCancelDelete)
		r.Post("/delete", routes.GolferDelete)
		r.Get("/export", routes.GolferExport)
		r.Get("/settings", routes.GolferSettings)
		r.Post("/settings", routes.GolferSettingsPost)
	})
	r.Route("/golfers/{name}", func(r chi.Router) {
		r.Use(middleware.GolferInfoHandler)
		r.Get("/", routes.GolferCheevos)
		r.Get("/holes", routes.GolferHoles)
		r.Get("/holes/{scoring}", routes.GolferHoles)
	})
	r.Get("/healthz", routes.Healthz)
	r.Get("/ideas", routes.Ideas)
	r.Get("/log-out", routes.LogOut)
	r.Get("/random", routes.Random)
	r.Get("/ng/random", routes.NGRandom)
	r.Route("/rankings", func(r chi.Router) {
		// Redirect some old URLs that got out.
		r.Get("/", redir("/rankings/holes/all/all/bytes"))
		r.Get("/holes", redir("/rankings/holes/all/all/bytes"))
		r.Get("/holes/all/all/all", redir("/rankings/holes/all/all/bytes"))
		r.Get("/langs/bytes", redir("/rankings/langs/all/bytes"))
		r.Get("/langs/chars", redir("/rankings/langs/all/chars"))
		r.Get("/medals", redir("/rankings/medals/all/all/all"))

		r.Get("/cheevos", routes.RankingsCheevos)
		r.Get("/cheevos/all", redir("/rankings/cheevos"))
		r.Get("/cheevos/{cheevo}", routes.RankingsCheevos)

		r.Get("/holes/{hole}/{lang}/{scoring}", routes.RankingsHoles)

		r.Get("/medals/{hole}/{lang}/{scoring}", routes.RankingsMedals)

		r.Get("/langs/{lang}/{scoring}", routes.RankingsLangs)
		r.Get("/solutions", routes.RankingsSolutions)
	})
	r.Get("/recent", routes.Recent)
	r.Get("/recent/{lang}", routes.Recent)
	r.Get("/scores/{hole}/{lang}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/all", routes.ScoresAll)
	r.Get("/scores/{hole}/{lang}/{scoring}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/{scoring}/{page}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/{scoring}/mini", routes.ScoresMini)
	r.Get("/sitemap.xml", routes.Sitemap)
	r.Post("/solution", routes.Solution)
	r.Get("/stats", routes.Stats)
	r.Get("/users/{name}", routes.User)

	certManager := autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"code.golf", "www.code.golf",
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
	if _, dev := os.LookupEnv("DEV"); dev {
		crt = "localhost.pem"
		key = "localhost-key.pem"
	} else {
		server.TLSConfig.GetCertificate = certManager.GetCertificate
	}

	// Every minute.
	go func() {
		for range time.Tick(time.Minute) {
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
				      SELECT user_id, 'its-over-9000'::cheevo
				        FROM points
				       WHERE points > 9000
				   UNION ALL
				      SELECT user_id, 'twenty-kiloleagues'::cheevo
				        FROM points
				       WHERE points >= 20000
				 ON CONFLICT DO NOTHING`,
			); err != nil {
				log.Println(err)
			}
		}
	}()

	// Every 5 minutes.
	go func() {
		// Various GitHub API requests.
		for range time.Tick(5 * time.Minute) {
			github.Run(db)
		}
	}()

	// Every hour.
	go func() {
		for range time.Tick(time.Hour) {
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

func redir(url string) http.HandlerFunc {
	return http.RedirectHandler(url, http.StatusPermanentRedirect).ServeHTTP
}
