package fractions

import (
	"math/rand"
	"strconv"
	"strings"
)

// a fraction (fract) is defined by its
// numerator and denominator in simplified
// form (n, d) and its scale factor (s)
type fract struct{ n, d, s int }

// most of this code is lifted from the
// intersection hole
func strconvunsimplifiedfrac(frac fract) (out string) {
	var outs []string
	outs = append(outs, strconv.Itoa(frac.n * frac.s))
	outs = append(outs, strconv.Itoa(frac.d * frac.s))
	return strings.Join(outs, "/")
}

func strconvsimplifiedfrac(frac fract) (out string) {
	var outs []string
	outs = append(outs, strconv.Itoa(frac.n))
	outs = append(outs, strconv.Itoa(frac.d))
	return strings.Join(outs, "/")
}

// if numerator and denominator are
// divisible by a number greater than
// 1 then the fraction is reducible
func isIrreducible(frac fract) bool {
	gcd := 2
    for gcd <= frac.d {
    	if frac.n % gcd == 0 && frac.d % gcd == 0 {
        	return false
        }
        gcd++
    }
    return true
}

// generator of random fraction with non-zero denominator
func fracGen() fract {
	return fract{
		n: rand.Intn(16),
		d: rand.Intn(15)+1,
		s: rand.Intn(15)+1,
	}
}

func fractions() (args []string, out string) {
	var outs []string

    /// some hardcoded cases
    f1 := fract{n: 1,  d: 1, s:1}
    f2 := fract{n: 1,  d: 1, s:10}
    f3 := fract{n: 1,  d: 2, s:1}
    f4 := fract{n: 1,  d: 2, s:10}
    f5 := fract{n: 2,  d: 1, s:1}
    f6 := fract{n: 2,  d: 1, s:10}
    f7 := fract{n: 0,  d: 1, s:1}
    f8 := fract{n: 15, d:14, s:15}
    
    // 1 over 1
    args = append(args, strconvunsimplifiedfrac(f1))
    outs = append(outs, strconvsimplifiedfrac(f1))
    
    // 10 over 10
    args = append(args, strconvunsimplifiedfrac(f2))
    outs = append(outs, strconvsimplifiedfrac(f2))
    
    // 1 over 2
    args = append(args, strconvunsimplifiedfrac(f3))
    outs = append(outs, strconvsimplifiedfrac(f3))
    
    // 10 over 20
    args = append(args, strconvunsimplifiedfrac(f4))
    outs = append(outs, strconvsimplifiedfrac(f4))
    
     // 2 over 1
    args = append(args, strconvunsimplifiedfrac(f5))
    outs = append(outs, strconvsimplifiedfrac(f5))
    
    // 20 over 10
    args = append(args, strconvunsimplifiedfrac(f6))
    outs = append(outs, strconvsimplifiedfrac(f6))
    
    // 0 over 1
    args = append(args, strconvunsimplifiedfrac(f7))
    outs = append(outs, strconvsimplifiedfrac(f7))
    
    // 225 over 210
    args = append(args, strconvunsimplifiedfrac(f8))
    outs = append(outs, strconvsimplifiedfrac(f8))

	//// generate 100 random cases
    cases := 0
    for cases < 100 {
    	f1 = fracGen()
        if isIrreducible(f1) {
        	args = append(args, strconvunsimplifiedfrac(f1))
    		outs = append(outs, strconvsimplifiedfrac(f1))
        	cases++
        }
    }

	// shuffle args and outputs in the same way
	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	out = strings.Join(outs, "\n")
	return
}
