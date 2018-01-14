package terminal

import "strconv"

var emptyStyle = style{}

type style struct {
	fgColor   uint8
	bgColor   uint8
	fgColorX  bool
	bgColorX  bool
	bold      bool
	faint     bool
	italic    bool
	underline bool
	strike    bool
	blink     bool
}

const (
	COLOR_NORMAL        = iota
	COLOR_GOT_38_NEED_5 = iota
	COLOR_GOT_48_NEED_5 = iota
	COLOR_GOT_38        = iota
	COLOR_GOT_48        = iota
)

// True if both styles are equal (or are the same object)
func (s *style) isEqual(o *style) bool {
	return s == o || *s == *o
}

// CSS classes that make up the style
func (s *style) asClasses() []string {
	var styles []string

	if s.fgColor > 0 && s.fgColor < 38 && !s.fgColorX {
		styles = append(styles, "term-fg"+strconv.Itoa(int(s.fgColor)))
	}
	if s.fgColor > 38 && !s.fgColorX {
		styles = append(styles, "term-fgi"+strconv.Itoa(int(s.fgColor)))

	}
	if s.fgColorX {
		styles = append(styles, "term-fgx"+strconv.Itoa(int(s.fgColor)))

	}

	if s.bgColor > 0 && s.bgColor < 48 && !s.bgColorX {
		styles = append(styles, "term-bg"+strconv.Itoa(int(s.bgColor)))
	}
	if s.bgColor > 48 && !s.bgColorX {
		styles = append(styles, "term-bgi"+strconv.Itoa(int(s.bgColor)))
	}
	if s.bgColorX {
		styles = append(styles, "term-bgx"+strconv.Itoa(int(s.bgColor)))
	}

	if s.bold {
		styles = append(styles, "term-fg1")
	}
	if s.faint {
		styles = append(styles, "term-fg2")
	}
	if s.italic {
		styles = append(styles, "term-fg3")
	}
	if s.underline {
		styles = append(styles, "term-fg4")
	}
	if s.blink {
		styles = append(styles, "term-fg5")
	}
	if s.strike {
		styles = append(styles, "term-fg9")
	}

	return styles
}

// True if style is empty
func (s *style) isEmpty() bool {
	return *s == style{}
}

// Add colours to an existing style, potentially returning
// a new style.
func (s *style) color(colors []string) *style {
	if len(colors) == 1 && (colors[0] == "0" || colors[0] == "") {
		// Shortcut for full style reset
		return &emptyStyle
	}

	newStyle := style(*s)
	s = &newStyle
	color_mode := COLOR_NORMAL

	for _, ccs := range colors {
		// If multiple colors are defined, i.e. \e[30;42m\e then loop through each
		// one, and assign it to s.fgColor or s.bgColor
		cc, err := strconv.ParseUint(ccs, 10, 8)
		if err != nil {
			continue
		}

		// State machine for XTerm colors, eg 38;5;150
		switch color_mode {
		case COLOR_GOT_38_NEED_5:
			if cc == 5 {
				color_mode = COLOR_GOT_38
			} else {
				color_mode = COLOR_NORMAL
			}
			continue
		case COLOR_GOT_48_NEED_5:
			if cc == 5 {
				color_mode = COLOR_GOT_48
			} else {
				color_mode = COLOR_NORMAL
			}
			continue
		case COLOR_GOT_38:
			s.fgColor = uint8(cc)
			s.fgColorX = true
			color_mode = COLOR_NORMAL
			continue
		case COLOR_GOT_48:
			s.bgColor = uint8(cc)
			s.bgColorX = true
			color_mode = COLOR_NORMAL
			continue
		}

		switch cc {
		case 0:
			// Reset all styles - don't use &emptyStyle here as we could end up adding colours
			// in this same action.
			s = &style{}
		case 1:
			s.bold = true
			s.faint = false
		case 2:
			s.faint = true
			s.bold = false
		case 3:
			s.italic = true
		case 4:
			s.underline = true
		case 5, 6:
			s.blink = true
		case 9:
			s.strike = true
		case 21, 22:
			s.bold = false
			s.faint = false
		case 23:
			s.italic = false
		case 24:
			s.underline = false
		case 25:
			s.blink = false
		case 29:
			s.strike = false
		case 38:
			color_mode = COLOR_GOT_38_NEED_5
		case 39:
			s.fgColor = 0
			s.fgColorX = false
		case 48:
			color_mode = COLOR_GOT_48_NEED_5
		case 49:
			s.bgColor = 0
			s.bgColorX = false
		case 30, 31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96, 97:
			s.fgColor = uint8(cc)
			s.fgColorX = false
		case 40, 41, 42, 43, 44, 45, 46, 47, 100, 101, 102, 103, 104, 105, 106, 107:
			s.bgColor = uint8(cc)
			s.bgColorX = false
		}
	}
	return s
}
