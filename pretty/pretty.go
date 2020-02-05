package pretty

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

const (
	day   = 24 * time.Hour
	week  = 7 * day
	month = 4 * week
)

// Comma returns a string of a int with thousand separators. Only 0 - 999,999.
func Comma(i int) string {
	if i > 999 {
		return fmt.Sprintf("%d,%03d", i/1000, i%1000)
	}

	return strconv.Itoa(i)
}

// Ordinal returns the ordinal of a int.
func Ordinal(i int) string {
	switch i % 10 {
	case 1:
		if i%100 != 11 {
			return "st"
		}
	case 2:
		if i%100 != 12 {
			return "nd"
		}
	case 3:
		if i%100 != 13 {
			return "rd"
		}
	}
	return "th"
}

// Time returns a fuzzy HTML <time> tag of a time.Time.
func Time(t time.Time) template.HTML {
	var sb strings.Builder

	rfc := t.Format(time.RFC3339)

	sb.WriteString("<time datetime=")
	sb.WriteString(rfc)
	sb.WriteString(" title=")
	sb.WriteString(rfc)
	sb.WriteRune('>')

	switch diff := time.Since(t); true {
	case diff < 2*time.Minute:
		sb.WriteString("1 min ago")
	case diff < 2*time.Hour:
		fmt.Fprintf(&sb, "%d mins ago", diff/time.Minute)
	case diff < 2*day:
		fmt.Fprintf(&sb, "%d hours ago", diff/time.Hour)
	case diff < 2*week:
		fmt.Fprintf(&sb, "%d days ago", diff/day)
	case diff < 2*month:
		fmt.Fprintf(&sb, "%d weeks ago", diff/week)
	default:
		fmt.Fprintf(&sb, "%d months ago", diff/month)
	}

	sb.WriteString("</time>")

	return template.HTML(sb.String())
}
