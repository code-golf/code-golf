package main

var intros = map[string]string{
	"99-bottles-of-beer": `<div class="beg hole"><a href=99-bottles-of-beer>99 Bottles of Beer<p>99 bottles of beer on the wall, 99 bottles of beer…</p></a><table>`,
	"fizz-buzz":          `<div class="beg hole"><a href=fizz-buzz>Fizz Buzz<p>Write a program that prints the numbers from 1 to 100…</p></a><table>`,
	"π":                  `<div class="beg hole"><a href=π>π<p>The ratio of a circle's circumference to its diameter…</p></a><table>`,
}

var preambles = map[string]string{
	"99-bottles-of-beer": "",
	"fizz-buzz": `<h1>Fizz Buzz</h1>

    <p>Write a program that prints the numbers from <b>1</b> to <b>100</b> (inclusive), each on their own line.</p>

    <p>If, however, the number is a multiple of <b>three</b> then print <b>Fizz</b> instead, and if the number is a multiple of <b>five</b> then print <b>Buzz</b>.</p>

    <p>For numbers which are multiples of <b>both three and five</b> then print <b>FizzBuzz</b>.</p>`,
	"π": `<h1>π</h1>

    <p>Print π (Pi) to the first 1,000 decimal places.</p>`,
}
