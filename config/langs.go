package config

import (
	"cmp"
	"encoding/hex"
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

var (
	// Standard languages.
	LangByID = map[string]*Lang{}
	LangList []*Lang

	// Experimental languages.
	ExpLangByID = map[string]*Lang{}
	ExpLangList []*Lang

	// All languages.
	AllLangByID = map[string]*Lang{}
	AllLangList []*Lang

	// Redirects.
	LangRedirects = map[string]string{}
)

type Lang struct {
	Args, Redirects, Env []string       `json:"-"`
	ArgsQuine            []string       `json:"-" toml:"args-quine"`
	Assembly             bool           `json:"-"`
	Digest               string         `json:"digest"`
	DigestTrunc          []byte         `json:"-"`
	Example              string         `json:"example"`
	Experiment           int            `json:"experiment,omitzero"`
	ID                   string         `json:"id"`
	LogoURL              string         `json:"logo-url"`
	Name                 string         `json:"name"`
	Released             toml.LocalDate `json:"released"`
	Size                 string         `json:"size"`
	Version              string         `json:"version"`
	Website              string         `json:"website"`
}

func initLangs() {
	// Digests.
	var digests map[string]string
	if !testing.Testing() {
		digestFile, err := os.Open("/lang-digests.json")
		if err != nil {
			panic(err)
		}
		defer digestFile.Close()

		if err := json.NewDecoder(digestFile).Decode(&digests); err != nil {
			panic(err)
		}
	}

	var langs map[string]*Lang
	unmarshal("data/langs.toml", &langs)

	for name, lang := range langs {
		lang.Example = strings.TrimSuffix(lang.Example, "\n")
		lang.ID = ID(name)
		lang.LogoURL = Assets["svg/"+lang.ID+".svg"]
		lang.Name = name

		// Digest & DigestTrunc (48-bit, 12 char trunc, like docker images).
		if lang.Digest = digests[lang.ID]; lang.Digest != "" {
			var err error
			if lang.DigestTrunc, err = hex.DecodeString(
				strings.TrimPrefix(lang.Digest, "sha256:")[:12],
			); err != nil {
				panic(err)
			}
		}

		// Redirects.
		for _, redirect := range lang.Redirects {
			LangRedirects[redirect] = lang.ID
		}

		AllLangByID[lang.ID] = lang
		AllLangList = append(AllLangList, lang)

		if lang.Experiment == 0 {
			LangByID[lang.ID] = lang
			LangList = append(LangList, lang)
		} else {
			ExpLangByID[lang.ID] = lang
			ExpLangList = append(ExpLangList, lang)
		}
	}

	// Case-insensitive sort.
	for _, langs := range [][]*Lang{LangList, ExpLangList, AllLangList} {
		slices.SortFunc(langs, func(a, b *Lang) int {
			return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		})
	}
}
