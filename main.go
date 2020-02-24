package main

import (
	"crypto/tls"
	"database/sql"
	"math/rand"
	"net/http"
	"syscall"
	"time"

	"github.com/code-golf/code-golf/github"
	"github.com/code-golf/code-golf/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	_, dev := syscall.Getenv("DEV")

	db, err := sql.Open("postgres", "")
	if err != nil {
		panic(err)
	}

	var r = chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Compress(5))
	r.Use(middleware.WithValue("db", db))

	host := "code-golf.io"
	if dev {
		host = "localhost"
	}
	redirDomain := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := "https://" + host + r.RequestURI
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	})

	// Redirect www (or any incorrect domain) to apex.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host == host {
				next.ServeHTTP(w, r)
			} else {
				redirDomain.ServeHTTP(w, r)
			}
		})
	})

	r.NotFound(routes.NotFound)

	r.Get("/", routes.Index)
	r.Get("/{hole}", routes.GETHole)
	r.Get("/about", routes.About)
	r.Get("/assets/{asset}", routes.Asset)
	r.Get("/callback", routes.Callback)
	r.Get("/favicon.ico", routes.Asset)
	r.Get("/feeds/{feed}", routes.Feed)
	r.Get("/ideas", routes.Ideas)
	r.Get("/log-out", routes.LogOut)
	r.Get("/random", routes.Random)
	r.Get("/recent", routes.Recent)
	r.Get("/robots.txt", routes.Robots)
	r.Get("/scores/{hole}/{lang}", routes.Scores)
	r.Get("/scores/{hole}/{lang}/{suffix}", routes.Scores)
	r.Post("/solution", routes.Solution)
	r.Get("/stats", routes.Stats)
	r.Get("/users/{user}", routes.User)

	certManager := autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("code-golf.io", "www.code-golf.io"),
	}

	server := &http.Server{
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
	if dev {
		crt = "localhost.pem"
		key = "localhost-key.pem"
	} else {
		server.TLSConfig.GetCertificate = certManager.GetCertificate
	}

	server.TLSConfig.BuildNameToCertificate()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)

		for {
			<-ticker.C
			github.Ideas(db)
			github.PullRequests(db)
			github.Stars(db)
		}
	}()

	println("Listening…")

	// Redirect HTTP to HTTPS and handle ACME challenges.
	go func() {
		panic(http.ListenAndServe(":80", certManager.HTTPHandler(redirDomain)))
	}()

	// Serve HTTPS.
	panic(server.ListenAndServeTLS(crt, key))
}
