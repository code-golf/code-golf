package hole

import (
	"fmt"
	"strings"
)

func presentNthChar(code string, nrLine int, pos int) string {
	line := strings.Split(code, "\n")[nrLine-1]
	if pos == 0 {
		return "/\n" + line + "\n"
	}
	return strings.Repeat(" ", pos-1) + "V\n" + line + "\n"
}

func palindromicQuineRequirements(code string) []struct {
	Name    string
	Pass    bool
	Message string
} {
	isCorrect := true
	var message string
	codeRune := []rune(code)
	mismatchAt := 0
	for i := 0; i < len(codeRune)/2; i++ {
		if codeRune[i] != codeRune[len(codeRune)-i-1] {
			isCorrect = false
			mismatchAt = i
			break
		}
	}
	if !isCorrect {
		currentLine := 1
		currentPos := 1
		for i, r := range codeRune {
			if r == '\n' {
				currentLine++
				currentPos = 0
			}
			if i == mismatchAt {
				message += fmt.Sprintf("the character %q at line %d position %d\n",
					codeRune[i], currentLine, currentPos)
				message += presentNthChar(code, currentLine, currentPos)
			}
			if i == len(code)-mismatchAt-1 {
				message += fmt.Sprintf("doesn't match the character %q at line %d position %d\n",
					codeRune[i], currentLine, currentPos)
				message += presentNthChar(code, currentLine, currentPos)
			}
			currentPos++
		}
	}
	return []struct {
		Name    string
		Pass    bool
		Message string
	}{
		{"code is palindromic", isCorrect, message},
	}
}
