package main

var preambles = map[string]string{
	"99-bottles-of-beer":       `<h1>99 Bottles of Beer</h1><p>Print the lyrics to the song 99 Bottles of Beer.</p>`,
	"emirp-numbers":            `<h1>Emirp Numbers</h1><p>An emirp (prime spelled backwards) is a prime number that results in a different prime when its decimal digits are reversed. For example both <b>13</b> and <b>31</b> are emirps.</p><p>Print all the emirps between <b>1</b> and <b>1000</b></p>`,
	"fibonacci":                `<h1>Fibonacci</h1><p>Print the first <b>31</b> Fibonacci numbers from <b>F<sub>0</sub> = 0</b> to <b>F<sub>30</sub> = 832040</b> (inclusive), each on a separate line.</p>`,
	"fizz-buzz":                `<h1>Fizz Buzz</h1><p>Print the numbers from <b>1</b> to <b>100</b> (inclusive), each on their own line.</p><p>If, however, the number is a multiple of <b>three</b> then print <b>Fizz</b> instead, and if the number is a multiple of <b>five</b> then print <b>Buzz</b>.</p><p>For numbers which are multiples of <b>both three and five</b> then print <b>FizzBuzz</b>.</p>`,
	"prime-numbers":            `<h1>Prime Numbers</h1><p>Print all the prime numbers between <b>1</b> and <b>100</b></p>`,
	"sierpiński-triangle":      `<h1>Sierpiński Triangle</h1><p>The Sierpiński triangle is a fractal and attractive fixed set with the overall shape of an equilateral triangle, subdivided recursively into smaller equilateral triangles.</p><p>A Sierpiński triangle of order 4 should look like this, print such an output:</p><pre>               ▲
              ▲ ▲
             ▲   ▲
            ▲ ▲ ▲ ▲
           ▲       ▲
          ▲ ▲     ▲ ▲
         ▲   ▲   ▲   ▲
        ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲
       ▲               ▲
      ▲ ▲             ▲ ▲
     ▲   ▲           ▲   ▲
    ▲ ▲ ▲ ▲         ▲ ▲ ▲ ▲
   ▲       ▲       ▲       ▲
  ▲ ▲     ▲ ▲     ▲ ▲     ▲ ▲
 ▲   ▲   ▲   ▲   ▲   ▲   ▲   ▲
▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲ ▲</pre>`,
	"arabic-to-roman-numerals": `<h1>Arabic to Roman Numerals</h1><p>For each arabic number argument print the same number in roman numerals.</p>`,
	"pascals-triangle":         `<h1>Pascal's Triangle</h1><p>Print the first <b>20 rows</b> of Pascal's triangle.</p>`,
	"seven-segment":            `<h1>Seven Segment</h1>`,
	"spelling-numbers":         `<h1>Spelling Numbers</h1><p>For each number argument print the same number spelled out in English.<p>The numbers will be in the range of <b>0</b> to <b>1,000</b> inclusive.</p>`,
	"e":                        `<h1>e</h1><p>Print e (Euler's number) to the first 1,000 decimal places.</p>`,
	"π":                        `<h1>π</h1><p>Print π (Pi) to the first 1,000 decimal places.</p>`,
}
