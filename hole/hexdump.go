package hole

func hexdump() []Run {
	return outputTestsWithSep("\n\n", shuffle(fixedTests("hexdump")))
}
