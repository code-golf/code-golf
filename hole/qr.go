package hole

import (
	"encoding/hex"
	"math/rand"
	"strings"

	"github.com/skip2/go-qrcode"
)

func randStr(len int) string {
	buf := make([]byte, len)
	for i := range buf {
		buf[i] = byte(33 + rand.Intn(94)) // randPrintableAscii: [33; 126]
	}
	return string(buf)
}

type matrix = [][]bool

func get(qr matrix, i, j int) byte {
	if qr[i][j] != (((i + j) & 1) == 0) {
		return 1
	}
	return 0
}

func get2(qr matrix, i, j int) byte {
	return 2*get(qr, i, j) + get(qr, i, j-1)
}

func get4(qr matrix, i, j, dir int) byte {
	return 4*get2(qr, i, j) + get2(qr, i+dir, j)
}

func get8(qr matrix, i, j, dir int) byte {
	return 16*get4(qr, i, j, dir) + get4(qr, i+2*dir, j, dir)
}

func getErrorCorrectionBlocks(qr matrix) []byte {
	return []byte{
		get8(qr, 9, 10, 1),
		get8(qr, 13, 10, 1),
		get8(qr, 17, 10, 1),
		get8(qr, 12, 8, -1),
		get8(qr, 9, 5, 1),
		get8(qr, 12, 3, -1),
		get8(qr, 9, 1, 1),
	}
}

func genQr(content string) matrix {
	qr, _ := qrcode.New(content, qrcode.Low)
	qr.DisableBorder = true
	return qr.Bitmap()
}

func qrIsStandard(qr matrix) bool {
	if !(qr[8][2] && !qr[8][3] && qr[8][4]) { // mask pattern: (i+j)%2 == 0
		return false
	}

	if get4(qr, 20, 20, -1) != 4 { // mode indicator: byte encoding
		return false
	}

	if get8(qr, 18, 20, -1) != 17 { // message length: 17
		return false
	}

	return true
}

func getStandardQr() (content string, qr matrix) {
	for {
		content = randStr(17)
		qr = genQr(content)

		if qrIsStandard(qr) {
			return
		}
	}
}

func qrToString(qr matrix, trimRight bool) string {
	var buf strings.Builder
	for i, row := range qr {
		if i > 0 {
			buf.WriteByte('\n')
		}
		if trimRight {
			lastTrue := -1
			for j, bit := range row {
				if bit {
					lastTrue = j
				}
			}
			row = row[:lastTrue+1]
		}
		for _, bit := range row {
			if bit {
				buf.WriteByte('#')
			} else {
				buf.WriteByte(' ')
			}
		}
	}
	return buf.String()
}

func qr(decoder bool) []Scorecard {
	content, qr := getStandardQr()
	qrString := qrToString(qr, !decoder)

	if decoder {
		return []Scorecard{{Args: []string{qrString}, Answer: content}}
	}

	return []Scorecard{{Args: []string{content}, Answer: qrString}}
}
