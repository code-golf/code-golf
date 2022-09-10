package pretty

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

// Bytes returns a string of integer bytes formatted as B/KiB/MiB.
func Bytes(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d.0 B", b)
	}
	div, exp := int(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KM"[exp])
}

// Comma returns a string of an integer with thousand separators.
func Comma(i int) string {
	switch {
	case i >= 1e6:
		return fmt.Sprintf("%d,%03d,%03d", i/1e6, i%1e6/1e3, i%1e3)
	case i >= 1e3:
		return fmt.Sprintf("%d,%03d", i/1e3, i%1e3)
	default:
		return strconv.Itoa(i)
	}
}

// Ordinal returns the ordinal of an integer.
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
//
//	   a min ago
//	  2 mins ago
//	         ...
//	 59 mins ago
//	 an hour ago
//	 2 hours ago
//	         ...
//	23 hours ago
//	   a day ago
//	  2 days ago
//	         ...
//	 28 days ago
//	  1 Jan 2020
//	         ...
//	 31 Dec 2020
func Time(t time.Time) template.HTML {
	const day = 24 * time.Hour

	var sb strings.Builder

	sb.WriteString(t.Format(
		"<time datetime=" + time.RFC3339Nano +
			` title="2 Jan 2006 15:04:05.000000 MST">`,
	))

	diff := time.Until(t)
	past := diff < 0
	if past {
		diff = -diff
	}

	if diff < 28*day {
		switch {
		case diff < 2*time.Minute:
			sb.WriteString("a min")
		case diff < time.Hour:
			fmt.Fprintf(&sb, "%d mins", diff/time.Minute)
		case diff < 2*time.Hour:
			sb.WriteString("an hour")
		case diff < day:
			fmt.Fprintf(&sb, "%d hours", diff/time.Hour)
		case diff < 2*day:
			sb.WriteString("a day")
		default:
			fmt.Fprintf(&sb, "%d days", diff/day)
		}

		if past {
			sb.WriteString(" ago")
		}
	} else {
		sb.WriteString(t.Format("2 Jan 2006"))
	}

	sb.WriteString("</time>")

	return template.HTML(sb.String())
}

func Title(s string) string {
	return strings.Title(s) //nolint:staticcheck
}
