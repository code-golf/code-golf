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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v4"
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

func main() {
	ctx := context.Background()

	testMakeCode()

	db, err := pgx.Connect(ctx, "user=code-golf sslmode=disable")
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

	tx, err := db.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback(ctx)

	userMap := map[string]int{}
	for _, info := range infoList {
		if _, ok := userMap[info.Login]; !ok {
			userMap[info.Login] = 1000 + len(userMap)
		}
	}

	for login, userID := range userMap {
		if _, err := tx.Exec(
			ctx,
			"INSERT INTO users (id, login) VALUES($1, $2)",
			userID, login,
		); err != nil {
			panic(err)
		}
	}

	for i, info := range infoList {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO solutions (  bytes,     chars,    code, hole, lang,
			                        scoring, submitted, user_id)
			                VALUES (     $1,        $2,      $3,   $4,   $5,
			                             $6,        $7,      $8)`,
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

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
}
