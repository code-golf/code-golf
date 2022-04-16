package hole

import "testing"

func BenchmarkArrows(b *testing.B)          { benchHole(b, arrows) }
func BenchmarkISBN(b *testing.B)            { benchHole(b, isbn) }
func BenchmarkOrdinalNumbers(b *testing.B)  { benchHole(b, ordinalNumbers) }
func BenchmarkSpellingNumbers(b *testing.B) { benchHole(b, spellingNumbers) }

func BenchmarkArabicToRoman(b *testing.B) {
	benchHole(b, func() ([]string, string) { return arabicToRoman(false) })
}

func BenchmarkRomanToArabic(b *testing.B) {
	benchHole(b, func() ([]string, string) { return arabicToRoman(true) })
}

func BenchmarkSudoku(b *testing.B) {
	benchHole(b, func() ([]string, string) { return sudoku(false) })
}

func BenchmarkSudokuV2(b *testing.B) {
	benchHole(b, func() ([]string, string) { return sudoku(true) })
}

func benchHole(b *testing.B, hole func() ([]string, string)) {
	for n := 0; n < b.N; n++ {
		hole()
	}
}
