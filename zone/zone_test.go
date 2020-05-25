package zone

import "testing"

func BenchmarkList(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, zone := range List() {
			_ = zone.String()
		}
	}
}
