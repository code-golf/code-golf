package main

// Update the local code-golf database with information from code.golf.
// Instead of storing the actual solution code, store 'a' * strokes.
// This is helpful for working on features that affect the appearance of the
// leaderboards, users page, recent page, and more.
// Run this script after running 'make dev' to start the server.
// Note that the database will be deleted when running 'make dev' again.
// To avoid this, restart the server directly using 'docker-compose up'.
// This script only works for a nearly empty database.
// Run 'docker-compose rm -f && docker-compose up' and then run this script.

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type solution struct {
	Hole, Lang, Login, Submitted string
	Strokes, Bytes               int
}

// Create a string containing the given number of characters and UTF-8 encoded bytes.
func makeCode(chars, bytes int) (result string) {
	delta := bytes - chars

	for _, replacement := range "😃景£" {
		replacementDelta := len(string(replacement)) - 1
		result += strings.Repeat(string(replacement), delta/replacementDelta)
		delta %= replacementDelta
	}

	result += strings.Repeat("a", chars-len([]rune(result)))
	return
}

func testMakeCode() {
	for x := 0; x < 10; x++ {
		for y := x; y <= 4*x; y++ {
			result := makeCode(x, y)
			if len([]rune(result)) != x {
				panic("unexpected rune count")
			}
			if len(result) != y {
				panic("unexpected byte count")
			}
		}
	}
}

func getInfo() (results []solution) {
	resp, err := http.Get("https://code.golf/scores/all-holes/all-langs/all")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		panic(err)
	}

	return
}

func main() {
	testMakeCode()

	db, err := sql.Open("postgres", "user=code-golf sslmode=disable")
	if err != nil {
		panic(err)
	}

	infoList := getInfo()
	userMap := map[string]int{}

	for _, info := range infoList {
		if _, ok := userMap[info.Login]; !ok {
			userMap[info.Login] = 1000 + len(userMap)
		}
	}

	for login, userID := range userMap {
		if _, err := db.Exec("INSERT INTO users (id, login) VALUES($1, $2)",
			userID, login,
		); err != nil {
			panic(err)
		}
	}

	for index, info := range infoList {
		code := makeCode(info.Strokes, info.Bytes)

		if _, err := db.Exec(
			"INSERT INTO code (code) VALUES ($1) ON CONFLICT DO NOTHING", code,
		); err != nil {
			panic(err)
		}

		if _, err := db.Exec(
			`INSERT INTO solutions (code_id, user_id, hole, lang, scoring)
			    SELECT id, $2, $3, $4, 'chars' FROM code WHERE code = $1`,
			code,
			userMap[info.Login],
			info.Hole,
			info.Lang,
		); err != nil {
			panic(err)
		}

		if index%1000 == 0 {
			fmt.Printf("Progress %d/%d\n", index+1, len(infoList))
		}
	}
}
