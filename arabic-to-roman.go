package main

var roman0 = []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
var roman1 = []string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
var roman2 = []string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
var roman3 = []string{"", "M", "MM", "MMM"}

func arabicToRoman(n int) string {
	return roman3[n%10000/1000] + roman2[n%1000/100] + roman1[n%100/10] + roman0[n%10]
}
