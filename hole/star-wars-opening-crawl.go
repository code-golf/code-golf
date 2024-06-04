package hole

func starWarsOpeningCrawl() []Run {
	return outputTestsWithSep("\n\n", shuffle(fixedTests("star-wars-opening-crawl")))
}
