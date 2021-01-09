package hole

import(
        "strconv"
)

func palindromicQuineRequirements(code string) []struct{Name    string
                                                        Pass    bool
							Message string}{
    isCorrect := true
    var message string
    codeRune := []rune(code)
    mismatchAt := 0
    for i := 0; i < len(codeRune)/2; i++ {
        if codeRune[i] != codeRune[len(codeRune)-i-1] {
            isCorrect = false
            mismatchAt = i;
        }
    }
    if !isCorrect {
        message = "the character " + string(codeRune[mismatchAt])
        message +=" at position " + strconv.Itoa(mismatchAt)
        message += "doesn't matches the character " + string(codeRune[len(codeRune) - mismatchAt - 1])
        message += " at position " + strconv.Itoa(len(codeRune) - mismatchAt - 1)
    }
    return []struct{Name string; Pass bool; Message string }{
			{"code is palindromic", isCorrect, message},
                   }
}
