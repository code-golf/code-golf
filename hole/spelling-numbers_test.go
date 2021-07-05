package hole

import "testing"

func BenchmarkSpellingNumbers(b *testing.B) {
	for n := 0; n < b.N; n++ {
		spellingNumbers()
	}
}
