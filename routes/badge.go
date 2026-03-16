package routes

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/code-golf/code-golf/golfer"
	"github.com/code-golf/code-golf/pretty"
	"github.com/code-golf/code-golf/session"
)

const (
	badgeHeight    = 20
	badgeRadius    = 3
	badgeCharWidth = 6
	badgePadding   = 10
	badgeFont      = "Verdana,Geneva,DejaVu Sans,sans-serif"
	badgeLabelFill = "#555"
)

type badgeMetric struct {
	Label string
	Color string
	Value func(info *golfer.GolferInfo) string
}

var badgeMetrics = map[string]badgeMetric{
	"holes": {
		Label: "holes",
		Color: "#e05d44",
		Value: func(info *golfer.GolferInfo) string {
			return fmt.Sprintf("%s / %s", pretty.Comma(info.Holes), pretty.Comma(info.HolesTotal))
		},
	},
	"langs": {
		Label: "langs",
		Color: "#4c1",
		Value: func(info *golfer.GolferInfo) string {
			return fmt.Sprintf("%s / %s", pretty.Comma(info.Langs), pretty.Comma(info.LangsTotal))
		},
	},
	"cheevos": {
		Label: "cheevos",
		Color: "#007ec6",
		Value: func(info *golfer.GolferInfo) string {
			return fmt.Sprintf("%s / %s", pretty.Comma(len(info.Cheevos)), pretty.Comma(info.CheevosTotal))
		},
	},
	"bytes": {
		Label: "bytes",
		Color: "#2a9d8f",
		Value: func(info *golfer.GolferInfo) string {
			return pretty.Comma(info.BytesPoints) + " pts"
		},
	},
	"chars": {
		Label: "chars",
		Color: "#7b5dd6",
		Value: func(info *golfer.GolferInfo) string {
			return pretty.Comma(info.CharsPoints) + " pts"
		},
	},
	"gold-medals": {
		Label: "gold medals",
		Color: "#dfb317",
		Value: func(info *golfer.GolferInfo) string {
			return pretty.Comma(info.Gold)
		},
	},
	"silver-medals": {
		Label: "silver medals",
		Color: "#9f9f9f",
		Value: func(info *golfer.GolferInfo) string {
			return pretty.Comma(info.Silver)
		},
	},
	"bronze-medals": {
		Label: "bronze medals",
		Color: "#b08d57",
		Value: func(info *golfer.GolferInfo) string {
			return pretty.Comma(info.Bronze)
		},
	},
	"diamonds": {
		Label: "diamonds",
		Color: "#4fc3f7",
		Value: func(info *golfer.GolferInfo) string {
			return pretty.Comma(info.Diamond)
		},
	},
	"unicorns": {
		Label: "unicorns",
		Color: "#c678dd",
		Value: func(info *golfer.GolferInfo) string {
			return pretty.Comma(info.Unicorn)
		},
	},
}

var badgeMetricAliases = map[string]string{
	"gold":    "gold-medals",
	"silver":  "silver-medals",
	"bronze":  "bronze-medals",
	"diamond": "diamonds",
	"unicorn": "unicorns",
}

// GET /golfers/{name}/badge/{metric}
func golferBadgeGET(w http.ResponseWriter, r *http.Request) {
	metric := strings.ToLower(param(r, "metric"))
	if before, ok := strings.CutSuffix(metric, ".svg"); ok {
		metric = before
	}
	if alias, ok := badgeMetricAliases[metric]; ok {
		metric = alias
	}

	metricDef, ok := badgeMetrics[metric]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	info := golfer.GetInfo(session.Database(r), param(r, "name"))
	if info == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	svg := badgeSVG(metricDef.Label, metricDef.Value(info), metricDef.Color)

	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "max-age=300, public")
	w.Write(svg)
}

func badgeSVG(label, value, color string) []byte {
	label = strings.TrimSpace(label)
	value = strings.TrimSpace(value)

	labelWidth := badgeTextWidth(label)
	valueWidth := badgeTextWidth(value)
	totalWidth := labelWidth + valueWidth

	labelTextLen := labelWidth - badgePadding
	valueTextLen := valueWidth - badgePadding

	labelX := labelWidth / 2
	valueX := labelWidth + (valueWidth / 2)

	aria := html.EscapeString(label + ": " + value)
	labelEsc := html.EscapeString(label)
	valueEsc := html.EscapeString(value)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" role="img" aria-label="%s">`, totalWidth, badgeHeight, aria)
	fmt.Fprintf(&sb, `<title>%s</title>`, aria)
	sb.WriteString(`<linearGradient id="s" x2="0" y2="100%">` +
		`<stop offset="0" stop-color="#fff" stop-opacity=".7"/>` +
		`<stop offset=".1" stop-color="#aaa" stop-opacity=".1"/>` +
		`<stop offset=".9" stop-color="#000" stop-opacity=".3"/>` +
		`<stop offset="1" stop-color="#000" stop-opacity=".5"/>` +
		`</linearGradient>`)
	fmt.Fprintf(&sb, `<mask id="r"><rect width="%d" height="%d" rx="%d" fill="#fff"/></mask>`, totalWidth, badgeHeight, badgeRadius)
	sb.WriteString(`<g mask="url(#r)">`)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="%s"/>`, labelWidth, badgeHeight, badgeLabelFill)
	fmt.Fprintf(&sb, `<rect x="%d" width="%d" height="%d" fill="%s"/>`, labelWidth, valueWidth, badgeHeight, color)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="url(#s)"/>`, totalWidth, badgeHeight)
	sb.WriteString(`</g>`)
	fmt.Fprintf(&sb, `<g fill="#fff" text-anchor="middle" font-family="%s" font-size="11">`, badgeFont)
	fmt.Fprintf(&sb, `<text x="%d" y="14" textLength="%d" lengthAdjust="spacingAndGlyphs">%s</text>`, labelX, labelTextLen, labelEsc)
	fmt.Fprintf(&sb, `<text x="%d" y="14" textLength="%d" lengthAdjust="spacingAndGlyphs">%s</text>`, valueX, valueTextLen, valueEsc)
	sb.WriteString(`</g></svg>`)

	return []byte(sb.String())
}

func badgeTextWidth(s string) int {
	return (badgeCharWidth * utf8.RuneCountInString(s)) + badgePadding
}
