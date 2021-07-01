package config

import (
	"os"
)

var (
	Host = "code.golf"
	Dev  bool
)

func init() {
	_, Dev = os.LookupEnv("DEV")
	if _, e2e := os.LookupEnv("E2E"); e2e {
		Host = "app:1443"
	} else if Dev {
		Host = "localhost"
	}
}
