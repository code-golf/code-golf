package config

import "encoding/json"

var Assets = map[string]string{}

func init() {
	if file, err := config.Open("data/assets.json"); err == nil {
		defer file.Close()

		var esbuild struct {
			Outputs map[string]struct{ EntryPoint string }
		}

		if err := json.NewDecoder(file).Decode(&esbuild); err != nil {
			panic(err)
		}

		for dist, src := range esbuild.Outputs {
			if src.EntryPoint != "" {
				Assets[src.EntryPoint] = "/" + dist
			}
		}
	}
}
