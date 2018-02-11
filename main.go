package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jraspass/code-golf/routes"
	"golang.org/x/crypto/acme/autocert"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Catch panics and turn them into 500s.
	defer func(start time.Time) {
		fmt.Printf(
			"%3dms %4s %s %s\n",
			time.Since(start).Nanoseconds()/1e6,
			r.Method,
			r.URL.Path,
			r.UserAgent(),
		)

		if r := recover(); r != nil {
			fmt.Fprint(os.Stderr, "<1>", r, "\n")
			http.Error(w, "500: It's Dead, Jim.", 500)
		}
	}(time.Now())

	switch r.Host {
	case "code-golf.io":
		w.Header()["Content-Language"] = []string{"en"}

		routes.Router.ServeHTTP(w, r)
	case "raspass.me":
		raspass(w, r)
	case "www.code-golf.io", "www.raspass.me":
		http.Redirect(w, r, "//"+r.Host[4:]+r.URL.String(), 301)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	certManager := autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"code-golf.io", "raspass.me", "www.code-golf.io", "www.raspass.me",
		),
	}

	server := &http.Server{
		Handler: &handler{},
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // HTTP/2-required.
			},
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			GetCertificate:           certManager.GetCertificate,
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	server.TLSConfig.BuildNameToCertificate()

	println("Listening…")

	// Redirect HTTP to HTTPS and handle ACME challenges.
	go func() { panic(http.ListenAndServe(":80", certManager.HTTPHandler(nil))) }()

	// Serve HTTPS.
	panic(server.ListenAndServeTLS("", ""))
}
