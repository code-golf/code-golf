package main

// Update the local code-golf database with information from code.golf.
// Instead of storing the actual solution code, a string with the same
// number of codepoints and UTF-8 bytes is generated and stored.
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
	"math/rand"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// Create a string containing the given number of characters and UTF-8 encoded bytes.
func makeCode(chars *int, bytes int) (result string) {
	if chars == nil {
		// Support assembly solutions without chars.
		return strings.Repeat("a", bytes)
	}

	delta := bytes - *chars

	for _, replacement := range "😃景£" {
		replacementDelta := len(string(replacement)) - 1
		result += strings.Repeat(string(replacement), delta/replacementDelta)
		delta %= replacementDelta
	}

	result += strings.Repeat("a", *chars-len([]rune(result)))
	return
}

func testMakeCode() {
	for x := 0; x < 10; x++ {
		for y := x; y <= 4*x; y++ {
			result := makeCode(&x, y)
			if len([]rune(result)) != x {
				panic("unexpected rune count")
			}
			if len(result) != y {
				panic("unexpected byte count")
			}
		}
	}
}

func getUserMap(db *sql.DB) (result map[string]int32) {
	result = map[string]int32{}

	rows, err := db.Query(`SELECT id, login FROM users`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int32
		var login string

		if err := rows.Scan(&id, &login); err != nil {
			panic(err)
		}

		result[login] = id
	}

	return
}

func getUnusedUserID(userMap map[string]int32) (result int32) {
	done := false
	for !done {
		done = true
		result = rand.Int31()
		for _, userID := range userMap {
			if result == userID {
				done = false
				break
			}
		}
	}

	return
}

func main() {
	testMakeCode()

	db, err := sql.Open("postgres", "user=code-golf sslmode=disable")
	if err != nil {
		panic(err)
	}

	res, err := http.Get("https://code.golf/scores/all-holes/all-langs/all")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var infoList []struct {
		Bytes                                 int
		Chars                                 *int
		Hole, Lang, Login, Scoring, Submitted string
	}
	if err := json.NewDecoder(res.Body).Decode(&infoList); err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	userMap := getUserMap(db)

	for _, info := range infoList {
		if _, ok := userMap[info.Login]; !ok {
			userID := getUnusedUserID(userMap)
			userMap[info.Login] = userID
			if _, err := tx.Exec("INSERT INTO users (id, login) VALUES($1, $2)",
				userID, info.Login,
			); err != nil {
				panic(err)
			}
		}
	}

	for i, info := range infoList {
		if _, err := tx.Exec(
			`INSERT INTO solutions (  bytes,     chars,    code, hole, lang,
			                        scoring, submitted, user_id)
			                VALUES (     $1,        $2,      $3,   $4,   $5,
			                             $6,        $7,      $8)
			ON CONFLICT            (user_id, hole, lang, scoring)
			  DO UPDATE SET bytes = excluded.bytes,
			                chars = excluded.chars,
			                 code = excluded.code,
			            submitted = excluded.submitted`,
			info.Bytes,
			info.Chars,
			makeCode(info.Chars, info.Bytes),
			info.Hole,
			info.Lang,
			info.Scoring,
			info.Submitted,
			userMap[info.Login],
		); err != nil {
			panic(err)
		}

		if i%1000 == 0 {
			fmt.Printf("Progress %d/%d\n", i+1, len(infoList))
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
