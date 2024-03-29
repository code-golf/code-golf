package config

import "slices"

type Option struct{ ID, Name string }

type Setting struct {
	Checkbox          bool
	Default, ID, Name string
	Options           []*Option
}

var Settings map[string][]*Setting

func init() {
	unmarshal("settings.toml", &Settings)

	for _, settings := range Settings {
		for _, setting := range settings {
			if len(setting.Options) == 0 {
				continue
			}

			// Default to the first option.
			if setting.Default == "" {
				setting.Default = setting.Options[0].ID
			}

			// A bit hacky, append all something options onto "All Something".
			switch setting.Options[0].Name {
			case "All Holes":
				for _, hole := range HoleList {
					setting.Options = append(setting.Options,
						&Option{ID: hole.ID, Name: hole.Name})
				}
			case "All Languages":
				for _, lang := range LangList {
					setting.Options = append(setting.Options,
						&Option{ID: lang.ID, Name: lang.Name})
				}
			}
		}
	}
}

func (s *Setting) ValueOrDefault(value string) string {
	// TODO Consider something more effecient like a hash?
	if slices.ContainsFunc(s.Options, func(o *Option) bool { return o.ID == value }) {
		return value
	}

	return s.Default
}
