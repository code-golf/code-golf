package routes

type Lang struct{ ID, Name string }
type Hole struct{ ID, Name, Difficulty, Preamble string }

var langs = []Lang{
	{"bash", "Bash"},
	{"haskell", "Haskell"},
	{"javascript", "JavaScript"},
	{"lisp", "Lisp"},
	{"lua", "Lua"},
	{"perl", "Perl"},
	{"perl6", "Perl 6"},
	{"php", "PHP"},
	{"python", "Python"},
	{"ruby", "Ruby"},
}

const morseTable = " International Morse Code.<ol><li>The length of a dot is one unit.<li>A dash is three units.<li>The space between parts of the same letter is one unit.<li>The space between letters is three units<li>The space between words is ten units</ol><table><tr><th>A<td>▄ ▄▄▄<tr><th>B<td>▄▄▄ ▄ ▄ ▄<tr><th>C<td>▄▄▄ ▄ ▄▄▄ ▄<tr><th>D<td>▄▄▄ ▄ ▄<tr><th>E<td>▄<tr><th>F<td>▄ ▄ ▄▄▄ ▄<tr><th>G<td>▄▄▄ ▄▄▄ ▄<tr><th>H<td>▄ ▄ ▄ ▄<tr><th>I<td>▄ ▄<tr><th>J<td>▄ ▄▄▄ ▄▄▄ ▄▄▄<tr><th>K<td>▄▄▄ ▄ ▄▄▄<tr><th>L<td>▄ ▄▄▄ ▄ ▄<tr><th>M<td>▄▄▄ ▄▄▄<tr><th>N<td>▄▄▄ ▄<tr><th>O<td>▄▄▄ ▄▄▄ ▄▄▄<tr><th>P<td>▄ ▄▄▄ ▄▄▄ ▄<tr><th>Q<td>▄▄▄ ▄▄▄ ▄ ▄▄▄<tr><th>R<td>▄ ▄▄▄ ▄<tr><th>S<td>▄ ▄ ▄<tr><th>T<td>▄▄▄<tr><th>U<td>▄ ▄ ▄▄▄<tr><th>V<td>▄ ▄ ▄ ▄▄▄<tr><th>W<td>▄ ▄▄▄ ▄▄▄<tr><th>X<td>▄▄▄ ▄ ▄ ▄▄▄<tr><th>Y<td>▄▄▄ ▄ ▄▄▄ ▄▄▄<tr><th>Z<td>▄▄▄ ▄▄▄ ▄ ▄<tr><th>0<td>▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄<tr><th>1<td>▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄<tr><th>2<td>▄ ▄ ▄▄▄ ▄▄▄ ▄▄▄<tr><th>3<td>▄ ▄ ▄ ▄▄▄ ▄▄▄<tr><th>4<td>▄ ▄ ▄ ▄ ▄▄▄<tr><th>5<td>▄ ▄ ▄ ▄ ▄<tr><th>6<td>▄▄▄ ▄ ▄ ▄ ▄<tr><th>7<td>▄▄▄ ▄▄▄ ▄ ▄ ▄<tr><th>8<td>▄▄▄ ▄▄▄ ▄▄▄ ▄ ▄<tr><th>9<td>▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄</table>"

var holes = []Hole{
	{
		"12-days-of-christmas", "12 Days of Christmas", "Medium",
		`Print the lyrics to the song <b>The 12 Days of Christmas</b>:</p><blockquote>On the First day of Christmas
My true love sent to me
A Partridge in a Pear Tree.

…

On the Twelfth day of Christmas
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
	}, {
		"99-bottles-of-beer", "99 Bottles of Beer", "Medium",
		"Print the lyrics to the song 99 Bottles of Beer.</p>",
	}, {
		"arabic-to-roman", "Arabic to Roman", "Slow",
		"For each arabic numeral argument print the same number in roman numerals.</p>",
	}, {
		"brainfuck", "Brainfuck", "Slow",
		"Brainfuck is a minimalistic esoteric programming language created by Urban Müller in 1993.<p>Assuming an infinitely large array, the entire Brainfuck alphabet matches the following pseudocode:<table><tr><th>&gt;<td>ptr++<tr><th>&lt;<td>ptr--<tr><th>+<td>array[ptr]++<tr><th>-<td>array[ptr]--<tr><th>.<td>print(chr(array[ptr]))<tr><th>[<td>while(array[ptr]){<tr><th>]<td>}</table><p>Write a program that will receive various Brainfuck programs as arguments and execute each program in turn.</p>",
	}, {
		"christmas-trees", "Christmas Trees", "Medium",
		`Print a size ascending range of Christmas trees using asterisks, ranging from size <b>3</b> to size <b>9</b>, each tree separated by a blank line.<p>A size <b>3</b> tree should look like this, with a single centered asterisk for the trunk:<pre>   *
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
	}, {
		"diamonds", "Diamonds", "Medium",
		`Print a size ascending range of Diamonds using the numbers <b>1</b> to <b>9</b>, ranging from size <b>1</b> to size <b>9</b>, each diamond separated by a blank line.<p>A size <b>1</b> diamond should look like this, a single centered <b>1</b>:<pre>         1</pre><p>With the largest size <b>9</b> diamond looking like this:<pre>         1
        121
       12321
      1234321
     123454321
    12345654321
   1234567654321
  123456787654321
 12345678987654321
  123456787654321
   1234567654321
    12345654321
     123454321
      1234321
       12321
        121
         1</pre>`,
	}, {
		"divisors", "Divisors", "Fast",
		"A number is a divisor of another number if it can divide into it with no remainder.<p>Print the positive divisors of each number from <b>1</b> to <b>100</b> inclusive, on their own line, with each divisor separated by a space.</p>",
	}, {
		"emirp-numbers", "Emirp Numbers", "Fast",
		"An emirp (prime spelled backwards) is a prime number that results in a different prime when its decimal digits are reversed. For example both <b>13</b> and <b>31</b> are emirps.</p><p>Print all the emirp numbers from <b>1</b> to <b>1000</b> inclusive, each on their own line.</p>",
	}, {
		"evil-numbers", "Evil Numbers", "Fast",
		"An evil number is a non-negative number that has an even number of 1s in its binary expansion.<p>Print all the evil numbers from <b>0</b> to <b>50</b> inclusive, each on their own line.<p>Numbers that are not evil are called <a href=odious-numbers>odious numbers</a>.</p>",
	}, {
		"fibonacci", "Fibonacci", "Fast",
		"Print the first <b>31</b> Fibonacci numbers from <b>F<sub>0</sub> = 0</b> to <b>F<sub>30</sub> = 832040</b> (inclusive), each on a separate line.</p>",
	}, {
		"fizz-buzz", "Fizz Buzz", "Fast",
		"Print the numbers from <b>1</b> to <b>100</b> inclusive, each on their own line.</p><p>If, however, the number is a multiple of <b>three</b> then print <b>Fizz</b> instead, and if the number is a multiple of <b>five</b> then print <b>Buzz</b>.</p><p>For numbers which are multiples of <b>both three and five</b> then print <b>FizzBuzz</b>.</p>",
	}, {
		"happy-numbers", "Happy Numbers", "Fast",
		"A happy number is defined by the following sequence: Starting with any positive integer, replace the number by the sum of the squares of its digits in base-ten, and repeat the process until the number either equals 1 (where it will stay), or it loops endlessly in a cycle that does not include 1. Those numbers for which this process ends in 1 are happy numbers, while those that do not end in 1 are sad numbers.<p>For example, 19 is happy, as the associated sequence is:</p><dl><dd>1<sup>2</sup> + 9<sup>2</sup> = 82<dd>8<sup>2</sup> + 2<sup>2</sup> = 68<dd>6<sup>2</sup> + 8<sup>2</sup> = 100<dd>1<sup>2</sup> + 0<sup>2</sup> + 0<sup>2</sup> = 1.</dl><p>Print all the happy numbers from <b>1</b> to <b>200</b> inclusive, each on their own line.</p>",
	}, {
		"morse-decoder", "Morse Decoder", "Medium",
		"Using ▄ (U+2584 Lower Half Block) to represent a dot, encode the argument from" + morseTable,
	}, {
		"morse-encoder", "Morse Encoder", "Medium",
		"Using ▄ (U+2584 Lower Half Block) to represent a dot, decode the argument into" + morseTable,
	}, {
		"niven-numbers", "Niven Numbers", "Fast",
		"A niven number is a non-negative number that is divisible by the sum of its digits.<p>Print all the niven numbers from <b>0</b> to <b>100</b> inclusive, each on their own line.</p>",
	}, {
		"odious-numbers", "Odious Numbers", "Fast",
		"An odious number is a non-negative number that has an odd number of 1s in its binary expansion.<p>Print all the odious numbers from <b>0</b> to <b>50</b> inclusive, each on their own line.<p>Numbers that are not odious are called <a href=evil-numbers>evil numbers</a>.</p>",
	}, {
		"pangram-grep", "Pangram Grep", "Medium",
		"A pangram is a sentence that uses every letter of a given alphabet.<p>Write a program that will receive various sentences as arguments and print those that are valid pangrams.</p>",
	}, {
		"pascals-triangle", "Pascal's Triangle", "Fast",
		"Print the first <b>20 rows</b> of Pascal's triangle.</p>",
	}, {
		"pernicious-numbers", "Pernicious Numbers", "Fast",
		"A pernicious number is a positive number where the sum of its binary expansion is a <a href=prime-numbers>prime number</a>.<p>For example, <b>5</b> is a pernicious number since <b>5 = 101<sub>2</sub></b> and <b>1 + 1 = 2</b>, which is prime.<p>Print all the pernicious numbers from <b>0</b> to <b>50</b> inclusive, each on their own line.</p>",
	}, {
		"poker", "Poker", "Slow",
		"Given various poker hands as arguments, print what type of hand each argument is.<p>The list of hands in ranking order are as follows:<table><tr><th>Royal Flush<td class=red>🃁🃎🃍🃋🃊<tr><th>Straight Flush<td>🃛🃚🃙🃘🃗<tr><th>Four of a Kind<td>🃕<span class=red>🃅🂵</span>🂥<span class=red>🃂</span><tr><th>Full House<td>🂦<span class=red>🂶🃆</span>🃞<span class=red>🂾</span><tr><th>Flush<td class=red>🃋🃉🃈🃄🃃<tr><th>Straight<td><span class=red>🃊</span>🂩<span class=red>🂸🃇</span>🃖<tr><th>Three of a Kind<td>🃝🂭<span class=red>🂽🂹</span>🂢<tr><th>Two Pair<td><span class=red>🂻</span>🂫🃓🂣<span class=red>🂲</span><tr><th>Pair<td>🂪<span class=red>🂺</span>🂨<span class=red>🂷</span>🃔<tr><th>High Card<td><span class=red>🃎🃍</span>🂧🂤<span class=red>🂳</span></table>",
	}, {
		"prime-numbers", "Prime Numbers", "Fast",
		"Print all the prime numbers from <b>1</b> to <b>100</b> inclusive, each on their own line.</p>",
	}, {
		"quine", "Quine", "Fast",
		"A <b>quine</b> is a non-empty computer program which takes no input and produces a copy of its own source code as its only output, produce such a program.<p>Trailing whitespace is <b>NOT</b> stripped from the output for this hole.</p>",
	}, {
		"roman-to-arabic", "Roman to Arabic", "Slow",
		"For each roman numeral argument print the same number in arabic numerals.</p>",
	}, {
		"rule-110", "Rule 110", "Slow",
		`Print the first <b>100</b> rows in the Rule 110 cellular automaton starting from an initial single living cell, which should begin like this:<pre>         █
        ██
       ███
      ██ █
     █████
    ██   █
   ███  ██
  ██ █ ███
 ███████ █
██     ███</pre><p>You can read more about on it and on its Turing-Completeness here: <a href="https://en.wikipedia.org/wiki/Rule_110">https://en.wikipedia.org/wiki/Rule_110</a>.`,
	}, {
		"seven-segment", "Seven Segment", "Medium",
		`Using pipes and underscores print the argument as if it were displayed on a seven segment display.<p>For example the number <b>0123456789</b> should be displayed as:<pre> _     _  _     _  _  _  _  _
| |  | _| _||_||_ |_   ||_||_|
|_|  ||_  _|  | _||_|  ||_| _|</pre>`,
	}, {
		"sierpiński-triangle", "Sierpiński Triangle", "Medium",
		`The Sierpiński triangle is a fractal and attractive fixed set with the overall shape of an equilateral triangle, subdivided recursively into smaller equilateral triangles.<p>A Sierpiński triangle of order 4 should look like this, print such an output:<pre>               ▲
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
	}, {
		"spelling-numbers", "Spelling Numbers", "Slow",
		"For each number argument print the same number spelled out in English.<p>The numbers will be in the range of <b>0</b> to <b>1,000</b> inclusive.</p>",
	}, {
		"λ", "λ", "Medium",
		"Print λ (Conway's constant) to the first 1,000 decimal places.</p>",
	}, {
		"π", "π", "Medium",
		"Print π (Pi) to the first 1,000 decimal places.</p>",
	}, {
		"τ", "τ", "Medium",
		"Print τ (Tau) to the first 1,000 decimal places.</p>",
	}, {
		"φ", "φ", "Medium",
		"Print φ (Phi) to the first 1,000 decimal places.</p>",
	}, {
		"𝑒", "𝑒", "Medium",
		"Print 𝑒 (Euler's number) to the first 1,000 decimal places.</p>",
	},
}

var langByID = map[string]Lang{}
var holeByID = map[string]Hole{}

func init() {
	for _, lang := range langs {
		langByID[lang.ID] = lang
	}

	for _, hole := range holes {
		holeByID[hole.ID] = hole
	}
}
