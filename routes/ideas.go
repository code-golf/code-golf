package routes

import "net/http"

// Ideas serves GET /ideas
func Ideas(w http.ResponseWriter, r *http.Request) {
	type idea struct {
		ID, ThumbsDown, ThumbsUp int
		Title                    string
	}

	rows, err := db.Query("SELECT * FROM ideas ORDER BY thumbs_up - thumbs_down DESC, title")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var ideas []idea

	for rows.Next() {
		var i idea

		if err := rows.Scan(&i.ID, &i.ThumbsDown, &i.ThumbsUp, &i.Title); err != nil {
			panic(err)
		}

		ideas = append(ideas, i)
	}

	render(w, r, http.StatusOK, "ideas", "Ideas", ideas)
}
