package hole

import (
	"math/rand"
	"strings"

	"github.com/skip2/go-qrcode"
)

func randStr(len int) string {
	buf := make([]byte, len)
	for i := range buf {
		buf[i] = byte(32 + rand.Intn(95)) // randPrintableAscii: [32; 126]
	}
	return string(buf)
}

func genQr(content string) [][]bool {
	qr, _ := qrcode.New(content, qrcode.Low)
	qr.DisableBorder = true
	return qr.Bitmap()
}

func qrIsStandard(b [][]bool) bool {
	if !(b[8][2] && !b[8][3] && b[8][4]) { // mask pattern: (i+j)%2 == 0
		return false
	}

	if !(b[20][20] && b[20][19] && !b[19][20] && b[19][19]) { // mode indicator: byte encoding
		return false
	}

	if !(!b[17][19] && !b[15][19]) { // message length: 17
		return false
	}

	return true
}

func qrToString(qr [][]bool) string {
	var buf strings.Builder
	for y := range qr {
		if y > 0 {
			buf.WriteByte('\n')
		}
		for x := range qr[y] {
			if qr[y][x] {
				buf.WriteRune('█')
			} else {
				buf.WriteByte(' ')
			}
		}
	}
	return buf.String()
}

func qr() ([]string, string) {
	for {
		content := randStr(17)
		qr := genQr(content)

		if qrIsStandard(qr) {
			qrString := qrToString(qr)
			return []string{qrString}, content
		}
	}
}
