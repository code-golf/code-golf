package hole

import (
	  "math/rand/v2"
	  "strings"
)

var fixedInputs = []string{
	  "lyndon", "factorization", "codegolf",
    "abcdefghijklmnopqrstuvwxyzyxwvutsrqponmlkjihgfedcbabcdefghijklmnopqrstuvwxyzyxwvutsrqponmlkjihgfedcba",
    "vdsxqiytqmptjqintoblcromtgjpujvhrwjoydwgbwezpfxwlvye",
    "rxexxecsudldbuosnzkvzvhlyofxczphartivkiehefdaffazlle",
    "dntkwgzapobnhyduatjqmkcfidgsdcnastyosbasaolycziwfphq",
    "aaaaaaaaaa",
}

var _ = answerFunc("continued-fractions", func() []Answer {
    const alphabet = "qwertzuiopasdfghjklyxcvbnm"
  
	  tests := make([]test, 100)

	  for i, input := range fixedInputs {
		    tests[i] = lyndonFactorizationTest(input)
	  }

	  for i := len(fixedInputs); i < len(tests); i++ {
        testLength := 1+rand.IntN(30)
        input := ""
        for range baseLength {
			      j := randInt(0, 25)
			      input += alphabet[j : j+1]
		    }
		    tests[i] = lyndonFactorizationTest(input)
	}
	return outputTests(shuffle(tests))
})

func lyndonFactorizationTest(input int) test {
	in := input
	
  var out strings.Builder
  i := 0
  
	for i < len(input) {
      j := i + 1
      k := i

      for j < len(input) && input[k] <= input[j] {
          if s[k] < s[j] { k = i } else { k += 1 }
          j += 1
      }
    
      for i <= k {
          out.writeString(input[i:i + j - k])
          out.writeByte(' ')
          i += j - k
      }
    
  }

	return test{in, strings.TrimSpace(out.String())}
}
