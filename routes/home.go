package routes

import "net/http"

// GET /
func homeGET(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Cards     []Card
		LangsUsed map[string]bool
	}{getHomeCards(r), map[string]bool{}}

	for _, card := range data.Cards {
		if card.Lang != nil {
			data.LangsUsed[card.Lang.ID] = true
		}
	}

	render(w, r, "home", data)
}
