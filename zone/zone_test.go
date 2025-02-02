package zone

import "testing"

func BenchmarkList(b *testing.B) {
	for b.Loop() {
		for _, zone := range List() {
			_ = zone.String()
		}
	}
}
