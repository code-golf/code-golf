package routes

func colour(rank int, points int, display string) string {
	if display == "points" {
		return colourPoints(points)
	}
	return colourRank(rank)
}

func colourRank(i int) string {
	if i <= 1 {
		return "yellow"
	}
	if i <= 2 {
		return "orange"
	}
	if i <= 3 {
		return "red"
	}
	if i <= 10 {
		return "purple"
	}
	if i <= 100 {
		return "blue"
	}
	return "green"
}

func colourPoints(i int) string {
	if i >= 950 {
		return "yellow"
	}
	if i >= 900 {
		return "orange"
	}
	if i >= 800 {
		return "red"
	}
	if i >= 700 {
		return "purple"
	}
	if i >= 500 {
		return "blue"
	}
	return "green"
}
