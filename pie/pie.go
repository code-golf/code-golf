package pie

import (
	"fmt"
	"html/template"
	"math"
	"strings"

	"github.com/code-golf/code-golf/pretty"
)

type Slice struct {
	Label    string
	Quantity int
	percent  float64
}

type Pie struct {
	slices   []Slice
	quantity float64
}

func New(slices []Slice) Pie {
	quantity := 0

	for _, slice := range slices {
		quantity += slice.Quantity
	}

	pie := Pie{slices, float64(quantity)}

	for i, slice := range slices {
		slices[i].percent = float64(slice.Quantity) / pie.quantity
	}

	return pie
}

func (pie Pie) HTML() template.HTML {
	var html strings.Builder

	html.WriteString(`<svg viewBox="-1 -1 2 2">`)

	agg, oldX, oldY := 0.0, 1.0, 0.0

	for _, slice := range pie.slices {
		largeArc := 0
		if slice.percent > .5 {
			largeArc = 1
		}

		agg += 2 * math.Pi * slice.percent
		x, y := math.Cos(agg), math.Sin(agg)

		fmt.Fprintf(
			&html,
			`<path d="M%v %vA1 1 0 %v 1 %v %vL0 0"/>`,
			oldX, oldY, largeArc, x, y,
		)

		oldX, oldY = x, y
	}

	html.WriteString("</svg><ul>")

	for _, slice := range pie.slices {
		fmt.Fprintf(
			&html,
			`<li class=color>%s<span>%v ≈ %2.1f%%</span>`,
			slice.Label,
			pretty.Comma(slice.Quantity),
			100*slice.percent,
		)
	}

	html.WriteString("</ul>")

	return template.HTML(html.String())
}
