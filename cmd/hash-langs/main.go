package main

import (
	"crypto/sha256"
	"encoding/json/v2"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const langDir = "/langs"

var h = sha256.New()

func main() {
	langs, err := os.ReadDir(langDir)
	if err != nil {
		panic(err)
	}

	digests := make(map[string]string, len(langs))

	for _, lang := range langs {
		h.Reset()

		dir := filepath.Join(langDir, lang.Name())
		if err := filepath.WalkDir(dir, walk); err != nil {
			panic(err)
		}

		digests[lang.Name()] = fmt.Sprintf("sha256:%x", h.Sum(nil))
	}

	file, err := os.Create("/lang-digests.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := json.MarshalWrite(file, digests); err != nil {
		panic(err)
	}
}

func walk(path string, d fs.DirEntry, err error) error {
	if err != nil || !d.Type().IsRegular() {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	return nil
}
