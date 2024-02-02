package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/buildkite/terminal-to-html/v3"
	"github.com/code-golf/code-golf/config"
	"github.com/code-golf/code-golf/hole"
)

// GET /sandbox
func sandboxGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Authors      []string
		HideDetails  bool
		Hole         *config.Hole
		Langs        map[string]*config.Lang
		RankingsView string
		Solutions    []map[string]string
	}{
		Langs: config.AllLangByID,
		Hole: &config.Hole{
			Name: "Sandbox",
			ID:   "sandbox",
		},
	}

	render(w, r, "hole-tabs", data, "Sandbox", "XYZ123")
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

	// TODO Should this be pushed lower?
	run.Stderr = string(terminal.Render([]byte(run.Stderr)))

	// The legacy single run we display, first failing or last overall.
	out := struct {
		// Legacy TitleCase attributes.
		Argv           []string
		Cheevos        []*config.Cheevo
		Err, Exp, Out  string
		ExitCode       int
		Pass, LoggedIn bool
		Took           time.Duration

		// Modern lowercase attributes.
		Runs []hole.Run `json:"runs"`
	}{
		Argv:     run.Args,
		Cheevos:  []*config.Cheevo{},
		Err:      run.Stderr,
		ExitCode: run.ExitCode,
		Exp:      run.Answer,
		LoggedIn: false,
		Out:      run.Stdout,
		Pass:     run.Pass,
		Runs:     []hole.Run{*run},
		Took:     run.Time,
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(&out); err != nil {
		panic(err)
	}
}
