package config

import "slices"

// TODO These structures are very similar to config/nav, consider sharing?
type Option struct {
	Group    bool
	ID, Name string
}

type Setting struct {
	Checkbox bool
	Default  any
	ID, Name string
	Options  []*Option
}

var Settings map[string][]*Setting

func initSettings() {
	unmarshal("data/settings.toml", &Settings)

	for _, settings := range Settings {
		for _, setting := range settings {
			// Simple boolean settings.
			if len(setting.Options) == 0 {
				// Default to false.
				if setting.Default == nil {
					setting.Default = false
				}

				continue
			}

			// Default to the first option.
			if setting.Default == nil {
				setting.Default = setting.Options[0].ID
			}

			// A bit hacky, append all something options onto "All Something".
			switch setting.Options[0].Name {
			case "All Holes":
				for _, hole := range HoleList {
					setting.Options = append(setting.Options,
						&Option{ID: hole.ID, Name: hole.Name})
				}

				setting.Options = append(setting.Options,
					&Option{Group: true, Name: "Experimental Holes"})

				for _, hole := range ExpHoleList {
					setting.Options = append(setting.Options,
						&Option{ID: hole.ID, Name: hole.Name})
				}
			case "All Languages":
				for _, lang := range LangList {
					setting.Options = append(setting.Options,
						&Option{ID: lang.ID, Name: lang.Name})
				}

				setting.Options = append(setting.Options,
					&Option{Group: true, Name: "Experimental Languages"})

				for _, lang := range ExpLangList {
					setting.Options = append(setting.Options,
						&Option{ID: lang.ID, Name: lang.Name})
				}
			}
		}
	}
}

func (s *Setting) FromFormValue(value string) any {
	// Simple boolean settings.
	if len(s.Options) == 0 {
		return value != ""
	}

	// TODO Consider something more effecient like a hash?
	if slices.ContainsFunc(s.Options, func(o *Option) bool {
		return o.ID == value && !o.Group
	}) {
		return value
	}

	return s.Default
}
