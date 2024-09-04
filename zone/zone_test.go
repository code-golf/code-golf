package zone

import "testing"

func BenchmarkList(b *testing.B) {
	for range b.N {
		for _, zone := range List() {
			_ = zone.String()
		}
	}
}
