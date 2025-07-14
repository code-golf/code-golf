package hole

import (
	"testing"

	"github.com/code-golf/code-golf/config"
)

func BenchmarkArabicToRoman(b *testing.B)   { benchHole(b, "arabic-to-roman") }
func BenchmarkArrows(b *testing.B)          { benchHole(b, "arrows") }
func BenchmarkISBN(b *testing.B)            { benchHole(b, "isbn") }
func BenchmarkOrdinalNumbers(b *testing.B)  { benchHole(b, "ordinal-numbers") }
func BenchmarkRomanToArabic(b *testing.B)   { benchHole(b, "roman-to-arabic") }
func BenchmarkSet(b *testing.B)             { benchHole(b, "set") }
func BenchmarkSIUnits(b *testing.B)         { benchHole(b, "si-units") }
func BenchmarkSpellingNumbers(b *testing.B) { benchHole(b, "spelling-numbers") }
func BenchmarkSudoku(b *testing.B)          { benchHole(b, "sudoku") }
func BenchmarkSudokuFillIn(b *testing.B)    { benchHole(b, "sudoku-fill-in") }

func benchHole(b *testing.B, id string) {
	answerFunc := config.AllHoleByID[id].AnswerFunc

	for b.Loop() {
		answerFunc()
	}
}
