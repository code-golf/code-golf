package assets

import (
	"embed"
	"encoding/json/v2"
	"io/fs"
	"log"
	"strings"
)

//go:embed all:out
var out embed.FS

var Files fs.FS
var Paths = map[string]string{}

func init() {
	var err error
	if Files, err = fs.Sub(out, "out"); err != nil {
		panic(err)
	}

	f, err := Files.Open("meta.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	var meta struct {
		Outputs map[string]struct {
			EntryPoint string `json:"entryPoint"`
		} `json:"outputs"`
	}

	if err := json.UnmarshalRead(f, &meta); err != nil {
		panic(err)
	}

	for out, src := range meta.Outputs {
		if in, _ := strings.CutPrefix(src.EntryPoint, "assets/"); in != "" {
			Paths[in] = "/" + strings.Replace(out, "out/", "", 1)
		}
	}
}

func Exists(name string) bool {
	f, _ := Files.Open(name)
	if f == nil {
		return false
	}

	defer f.Close()
	return true
}
