package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/JRaspass/code-golf/routes"
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

	routes.Router.ServeHTTP(w, r)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	certManager := autocert.Manager{
		Cache:  autocert.DirCache("certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"code-golf.io", "www.code-golf.io",
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
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	_, dev := syscall.Getenv("DEV")

	if !dev {
		server.TLSConfig.GetCertificate = certManager.GetCertificate
	}

	server.TLSConfig.BuildNameToCertificate()

	go func() {
		ticker := time.NewTicker(time.Minute)

		for {
			<-ticker.C
			routes.Stars()
		}
	}()

	println("Listening…")

	// Redirect HTTP to HTTPS and handle ACME challenges.
	go func() { panic(http.ListenAndServe(":80", certManager.HTTPHandler(nil))) }()

	// Serve HTTPS.
	if dev {
		panic(server.ListenAndServeTLS("localhost.pem", "localhost-key.pem"))
	} else {
		panic(server.ListenAndServeTLS("", ""))
	}
}
