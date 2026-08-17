package hole

func isValidChessMove(piece, f1, r1, f2, r2 byte) bool {
	dx := int(f1) - int(f2)
	if dx < 0 {
		dx = -dx
	}

	dy := int(r1) - int(r2)
	if dy < 0 {
		dy = -dy
	}

	switch piece {
	case 'K':
		return dx <= 1 && dy <= 1
	case 'Q':
		return dx == 0 || dy == 0 || dx == dy
	case 'R':
		return dx == 0 || dy == 0
	case 'B':
		return dx == dy
	case 'N':
		return (dx == 1 && dy == 2) || (dx == 2 && dy == 1)
	}

	return false
}

var _ = answerFunc("chess-move-validation", func() []Answer {
	return outputTests(
		chessMovesTests(128),
		chessMovesTests(256),
		chessMovesTests(512),
	)
})

func chessMovesTests(n int) []test {
	pieces := []byte{'K', 'Q', 'R', 'B', 'N'}
	var validArgs, invalidArgs []string

	for _, p := range pieces {
		for f1 := byte('a'); f1 <= 'h'; f1++ {
			for r1 := byte('1'); r1 <= '8'; r1++ {
				for f2 := byte('a'); f2 <= 'h'; f2++ {
					for r2 := byte('1'); r2 <= '8'; r2++ {
						if f1 == f2 && r1 == r2 {
							continue // Must move to a new square
						}

						arg := string([]byte{p, f1, r1, f2, r2})
						if isValidChessMove(p, f1, r1, f2, r2) {
							validArgs = append(validArgs, arg)
						} else {
							invalidArgs = append(invalidArgs, arg)
						}
					}
				}
			}
		}
	}

	args := append(shuffle(validArgs)[:n], shuffle(invalidArgs)[:n]...)
	tests := make([]test, len(args))

	for i, arg := range shuffle(args) {
		tests[i].in = arg

		if isValidChessMove(arg[0], arg[1], arg[2], arg[3], arg[4]) {
			tests[i].out = arg
		}
	}

	return tests
}
