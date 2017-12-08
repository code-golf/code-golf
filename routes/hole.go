package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var preambles = map[string]string{
	"12-days-of-christmas": `<h1>12 Days of Christmas</h1><p>Print the lyrics to the song <b>The 12 Days of Christmas</b>:</p><blockquote>On the first day of Christmas
My true love sent to me
A Partridge in a Pear Tree.

…

On the twelfth day of Christmas
My true love sent to me
Twelve Drummers Drumming,
Eleven Pipers Piping,
Ten Lords-a-Leaping,
Nine Ladies Dancing,
Eight Maids-a-Milking,
Seven Swans-a-Swimming,
Six Geese-a-Laying,
Five Gold Rings,
Four Calling Birds,
Three French Hens,
Two Turtle Doves, and
A Partridge in a Pear Tree.</blockquote>`,
	"99-bottles-of-beer": `<h1>99 Bottles of Beer</h1><p>Print the lyrics to the song 99 Bottles of Beer.</p>`,
	"arabic-to-roman":    `<h1>Arabic to Roman</h1><p>For each arabic numeral argument print the same number in roman numerals.</p>`,
	"christmas-trees": `<h1>Christmas Trees</h1><p>Print a size ascending range of Christmas trees using asterisks, ranging from size <b>3</b> to size <b>9</b>, each tree separated by a blank line.<p>A size <b>3</b> tree should look like this, with a single centered asterisk for the trunk:<pre>   *
  ***
 *****
   *</pre><p>With the largest size <b>9</b> tree looking like this:<pre>         *
        ***
       *****
      *******
     *********
    ***********
   *************
  ***************
 *****************
         *</pre>`,
	"emirp-numbers":      `<h1>Emirp Numbers</h1><p>An emirp (prime spelled backwards) is a prime number that results in a different prime when its decimal digits are reversed. For example both <b>13</b> and <b>31</b> are emirps.</p><p>Print all the emirp numbers from <b>1</b> to <b>1000</b> inclusive, each on their own line.</p>`,
	"evil-numbers":       `<h1>Evil Numbers</h1><p>An evil number is a non-negative number that has an even number of 1s in its binary expansion.<p>Print all the evil numbers from <b>0</b> to <b>50</b> inclusive, each on their own line.<p>Numbers that are not evil are called <a href=odious-numbers>odious numbers</a>.</p>`,
	"fibonacci":          `<h1>Fibonacci</h1><p>Print the first <b>31</b> Fibonacci numbers from <b>F<sub>0</sub> = 0</b> to <b>F<sub>30</sub> = 832040</b> (inclusive), each on a separate line.</p>`,
	"fizz-buzz":          `<h1>Fizz Buzz</h1><p>Print the numbers from <b>1</b> to <b>100</b> inclusive, each on their own line.</p><p>If, however, the number is a multiple of <b>three</b> then print <b>Fizz</b> instead, and if the number is a multiple of <b>five</b> then print <b>Buzz</b>.</p><p>For numbers which are multiples of <b>both three and five</b> then print <b>FizzBuzz</b>.</p>`,
	"happy-numbers":      `<h1>Happy Numbers</h1><p>A happy number is defined by the following sequence: Starting with any positive integer, replace the number by the sum of the squares of its digits in base-ten, and repeat the process until the number either equals 1 (where it will stay), or it loops endlessly in a cycle that does not include 1. Those numbers for which this process ends in 1 are happy numbers, while those that do not end in 1 are sad numbers.<p>For example, 19 is happy, as the associated sequence is:</p><dl><dd>1<sup>2</sup> + 9<sup>2</sup> = 82<dd>8<sup>2</sup> + 2<sup>2</sup> = 68<dd>6<sup>2</sup> + 8<sup>2</sup> = 100<dd>1<sup>2</sup> + 0<sup>2</sup> + 0<sup>2</sup> = 1.</dl><p>Print all the happy numbers from <b>1</b> to <b>200</b> inclusive, each on their own line.</p>`,
	"odious-numbers":     `<h1>Odious Numbers</h1><p>An odious number is a non-negative number that has an odd number of 1s in its binary expansion.<p>Print all the odious numbers from <b>0</b> to <b>50</b> inclusive, each on their own line.<p>Numbers that are not odious are called <a href=evil-numbers>evil numbers</a>.</p>`,
	"pangram-grep":       `<h1>Pangram Grep</h1><p>A pangram is a sentence that uses every letter of a given alphabet.<p>Write a program that will receive various sentences as arguments and print those that are valid pangrams.</p>`,
	"pascals-triangle":   `<h1>Pascal's Triangle</h1><p>Print the first <b>20 rows</b> of Pascal's triangle.</p>`,
	"pernicious-numbers": `<h1>Pernicious Numbers</h1><p>A pernicious number is a positive number where the sum of its binary expansion is a <a href=prime-numbers>prime number</a>.<p>For example, <b>5</b> is a pernicious number since <b>5 = 101<sub>2</sub></b> and <b>1 + 1 = 2</b>, which is prime.<p>Print all the pernicious numbers from <b>0</b> to <b>50</b> inclusive, each on their own line.</p>`,
	"prime-numbers":      `<h1>Prime Numbers</h1><p>Print all the prime numbers from <b>1</b> to <b>100</b> inclusive, each on their own line.</p>`,
	"quine":              `<h1>Quine</h1><p>A <b>quine</b> is a non-empty computer program which takes no input and produces a copy of its own source code as its only output, produce such a program.<p>Trailing whitespace is <b>NOT</b> stripped from the output for this hole.</p>`,
	"roman-to-arabic":    `<h1>Roman to Arabic</h1><p>For each roman numeral argument print the same number in arabic numerals.</p>`,
	"seven-segment": `<h1>Seven Segment</h1><p>Using pipes and underscores print the argument as if it were displayed on a seven segment display.<p>For example the number <b>0123456789</b> should be displayed as:<pre> _     _  _     _  _  _  _  _
| |  | _| _||_||_ |_   ||_||_|
|_|  ||_  _|  | _||_|  ||_| _|</pre>`,
	"spelling-numbers": `<h1>Spelling Numbers</h1><p>For each number argument print the same number spelled out in English.<p>The numbers will be in the range of <b>0</b> to <b>1,000</b> inclusive.</p>`,
	"sierpiński-triangle": `<h1>Sierpiński Triangle</h1><p>The Sierpiński triangle is a fractal and attractive fixed set with the overall shape of an equilateral triangle, subdivided recursively into smaller equilateral triangles.<p>A Sierpiński triangle of order 4 should look like this, print such an output:<pre>               ▲
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
	"π": `<h1>π</h1><p>Print π (Pi) to the first 1,000 decimal places.</p>`,
	"φ": `<h1>φ</h1><p>Print φ (Phi) to the first 1,000 decimal places.</p>`,
	"𝑒": `<h1>e</h1><p>Print 𝑒 (Euler's number) to the first 1,000 decimal places.</p>`,
	"τ": `<h1>τ</h1><p>Print τ (Tau) to the first 1,000 decimal places.</p>`,
}

func hole(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	hole := r.URL.Path[1:]

	userID := printHeader(w, r, 200)

	w.Write([]byte(
		"<script async src=" + holeJsPath + "></script><div id=status><div>" +
			"<h2>Program Arguments</h2><pre id=Arg></pre>" +
			"<h2>Standard Error</h2><pre id=Err></pre>" +
			"<h2>Expected Output</h2><pre id=Exp></pre>" +
			"<h2>Standard Output</h2><pre id=Out></pre>" +
			"</div></div><main",
	))

	if userID == 0 {
		w.Write([]byte(
			"><div id=alert>Please " +
				`<a href="//github.com/login/oauth/authorize?` +
				`client_id=7f6709819023e9215205&scope=user:email">` +
				"Login with GitHub</a> in order to save solutions " +
				"and appear on the leaderboards.</div",
		))
	} else {
		var html []byte

		// Fetch the latest successful lang.
		if err := db.QueryRow(
			`SELECT CONCAT(' data-lang=', lang)
			   FROM solutions
			  WHERE user_id = $1 AND hole = $2`,
			userID, hole,
		).Scan(&html); err == nil {
			w.Write(html)
		} else if err != sql.ErrNoRows {
			panic(err)
		}

		// Fetch all the code per lang.
		if err := db.QueryRow(
			`SELECT STRING_AGG(CONCAT(
			            ' data-',
			            lang,
			            '="',
			            REPLACE(code, '"', '&#34;'),
			            '"'
			        ), '')
			   FROM solutions
			  WHERE user_id = $1 AND hole = $2`,
			userID, hole,
		).Scan(&html); err == nil {
			w.Write(html)
		} else if err != sql.ErrNoRows {
			panic(err)
		}
	}

	w.Write([]byte(
		">" + preambles[hole] +
			"<a class=tab href=#bash>Bash<span></span></a>" +
			"<a class=tab href=#javascript>JS<span></span></a>" +
			"<a class=tab href=#perl>Perl<span></span></a>" +
			"<a class=tab href=#perl6>Perl 6<span></span></a>" +
			"<a class=tab href=#php>PHP<span></span></a>" +
			"<a class=tab href=#python>Python<span></span></a>" +
			"<a class=tab href=#ruby>Ruby<span></span></a>" +
			"<input type=submit value=Run>",
	))
}
