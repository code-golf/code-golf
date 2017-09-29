package main

var intros = map[string]string{
	"99-bottles-of-beer":       `<div class="beg hole"><a href=99-bottles-of-beer>99 Bottles of Beer<p>99 bottles of beer on the wall, 99 bottles of beer…</p></a><table>`,
	"arabic-to-roman-numerals": `<div class="int hole"><a href=arabic-to-roman-numerals>Arabic to Roman<p>Convert Hindu–Arabic numerals to Roman numerals…</p></a><table>`,
	"fibonacci":                `<div class="beg hole"><a href=fibonacci>Fibonacci<p>Each number is the sum of the two preceding numbers…</p></a><table>`,
	"fizz-buzz":                `<div class="beg hole"><a href=fizz-buzz>Fizz Buzz<p>Write a program that prints the numbers from 1 to 100…</p></a><table>`,
	"pascals-triangle":         `<div class="int hole"><a href=pascals-triangle>Pascal's Triangle<p>Blaise Pascal's arithmetic and geometric figure…</p></a><table>`,
	"seven-segment":            `<div class="int hole"><a href=seven-segment>Seven Segment<p>Using pipes and underscores print a seven segment display…</p></a><table>`,
	"spelling-numbers":         `<div class="int hole"><a href=spelling-numbers>Spelling Numbers</a><table>`,
	"e":                        `<div class="adv hole"><a href=e>e<p>The unique number whose natural logarithm is equal to one…</p></a><table>`,
	"π":                        `<div class="adv hole"><a href=π>π<p>The ratio of a circle's circumference to its diameter…</p></a><table>`,
}

var preambles = map[string]string{
	"99-bottles-of-beer":       `<h1>99 Bottles of Beer</h1><p>Print the lyrics to the song 99 Bottles of Beer.</p>`,
	"arabic-to-roman-numerals": `<h1>Arabic to Roman Numerals</h1><p>For each arabic number argument print the same number in roman numerals.</p>`,
	"fibonacci":                `<h1>Fibonacci</h1><p>Print the first <b>31</b> Fibonacci numbers from <b>F<sub>0</sub> = 0</b> to <b>F<sub>30</sub> = 832040</b> (inclusive), each on a separate line.</p>`,
	"fizz-buzz":                `<h1>Fizz Buzz</h1><p>Print the numbers from <b>1</b> to <b>100</b> (inclusive), each on their own line.</p><p>If, however, the number is a multiple of <b>three</b> then print <b>Fizz</b> instead, and if the number is a multiple of <b>five</b> then print <b>Buzz</b>.</p><p>For numbers which are multiples of <b>both three and five</b> then print <b>FizzBuzz</b>.</p>`,
	"pascals-triangle":         `<h1>Pascal's Triangle</h1><p>Print the first <b>20 rows</b> of Pascal's triangle.</p>`,
	"seven-segment":            `<h1>Seven Segment</h1>`,
	"spelling-numbers":         `<h1>Spelling Numbers</h1><p>For each number argument print the same number spelled out in English.<p>The numbers will be in the range of <b>0</b> to <b>1,000</b> inclusive.</p>`,
	"e":                        `<h1>e</h1><p>Print e (Euler's number) to the first 1,000 decimal places.</p>`,
	"π":                        `<h1>π</h1><p>Print π (Pi) to the first 1,000 decimal places.</p>`,
}
