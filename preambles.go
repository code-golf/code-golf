package main

var intros = map[string]string{
	"99-bottles-of-beer":       `<div class="beg hole"><a href=99-bottles-of-beer>99 Bottles of Beer<p>99 bottles of beer on the wall, 99 bottles of beer…</p></a><table>`,
	"arabic-to-roman-numerals": `<div class="int hole"><a href=arabic-to-roman-numerals>Arabic to Roman Numerals</a><table>`,
	"fizz-buzz":                `<div class="beg hole"><a href=fizz-buzz>Fizz Buzz<p>Write a program that prints the numbers from 1 to 100…</p></a><table>`,
	"pascals-triangle":         `<div class="int hole"><a href=pascals-triangle>Pascal's Triangle<p>Blaise Pascal's arithmetic and geometric figure…</p></a><table>`,
	"π":                        `<div class="adv hole"><a href=π>π<p>The ratio of a circle's circumference to its diameter…</p></a><table>`,
}

var preambles = map[string]string{
	"99-bottles-of-beer": `<h1>99 Bottles of Beer</h1>

    <p>Print the lyrics to the song 99 Bottles of Beer.</p>`,
	"arabic-to-roman-numerals": `<h1>Arabic to Roman Numerals</h1>`,
	"fizz-buzz": `<h1>Fizz Buzz</h1>

    <p>Print the numbers from <b>1</b> to <b>100</b> (inclusive), each on their own line.</p>

    <p>If, however, the number is a multiple of <b>three</b> then print <b>Fizz</b> instead, and if the number is a multiple of <b>five</b> then print <b>Buzz</b>.</p>

    <p>For numbers which are multiples of <b>both three and five</b> then print <b>FizzBuzz</b>.</p>`,
	"pascals-triangle": `<h1>Pascal's Triangle</h1>

    <p>Prints the first <b>20 rows</b> of Pascal's triangle.</p>`,
	"π": `<h1>π</h1>

    <p>Print π (Pi) to the first 1,000 decimal places.</p>`,
}
