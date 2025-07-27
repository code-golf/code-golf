package routes

import (
	"encoding/json"
	"net/http"

	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/hole"
)

// GET /sandbox
func sandboxGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Authors                  []string
		HideDetails              bool
		Hole, PrevHole, NextHole *config.Hole
		Langs                    map[string]*config.Lang
		RankingsView             string
		Solutions                []map[string]string
		IsSponsor                bool
		HasNotes                 bool
	}{
		Langs: config.AllLangByID,
		Hole: &config.Hole{
			Name: "Sandbox",
			ID:   "sandbox",
		},
	}

	render(w, r, "hole-tabs", data, "Sandbox", "Sandbox")
}

// POST /sandbox
func sandboxPOST(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Code string
		Lang string
		Args []string
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		panic(err)
	}
	defer r.Body.Close()

	langObj := config.AllLangByID[in.Lang]
	if langObj == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 128 KiB, >= because arguments needs a null termination.
	if len(in.Code) >= 128*1024 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	run := hole.PlaySandbox(r.Context(), langObj, in.Code, in.Args)

	out := struct {
		Runs []hole.Run `json:"runs"`
	}{
		Runs: []hole.Run{*run},
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(&out); err != nil {
		panic(err)
	}
}
