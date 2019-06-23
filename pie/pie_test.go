package pie

import "testing"

func BenchmarkHTML(b *testing.B) {
	pie := New([]Slice{
		{Label: "Art", Quantity: 6},
		{Label: "Computing", Quantity: 3},
		{Label: "Gaming", Quantity: 4},
		{Label: "Mathematics", Quantity: 6},
		{Label: "Sequence", Quantity: 12},
		{Label: "Transform", Quantity: 7},
	})

	for n := 0; n < b.N; n++ {
		pie.HTML()
	}
}
