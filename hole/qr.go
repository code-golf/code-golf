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

type matrix [][]bool

func (qr matrix) get(i, j int) byte {
	if qr[i][j] != (((i + j) & 1) == 0) {
		return 1
	}
	return 0
}

func (qr matrix) get2(i, j int) byte {
	return 2*qr.get(i, j) + qr.get(i, j-1)
}

func (qr matrix) get4(i, j, dir int) byte {
	return 4*qr.get2(i, j) + qr.get2(i+dir, j)
}

func (qr matrix) get8(i, j, dir int) byte {
	return 16*qr.get4(i, j, dir) + qr.get4(i+2*dir, j, dir)
}

func (qr matrix) getErrorCorrectionBlocks() []byte {
	return []byte{
		qr.get8(9, 10, 1),
		qr.get8(13, 10, 1),
		qr.get8(17, 10, 1),
		qr.get8(12, 8, -1),
		qr.get8(9, 5, 1),
		qr.get8(12, 3, -1),
		qr.get8(9, 1, 1),
	}
}

func genQr(content string) matrix {
	qr, _ := qrcode.New(content, qrcode.Low)
	qr.DisableBorder = true
	return qr.Bitmap()
}

func (qr matrix) isStandard() bool {
	if !(qr[8][2] && !qr[8][3] && qr[8][4]) { // mask pattern: (i+j)%2 == 0
		return false
	}

	if qr.get4(20, 20, -1) != 4 { // mode indicator: byte encoding
		return false
	}

	if qr.get8(18, 20, -1) != 17 { // message length: 17
		return false
	}

	return true
}

func getStandardQr() (content string, qr matrix) {
	for {
		content = randStr(17)
		qr = genQr(content)

		if qr.isStandard() {
			return
		}
	}
}

func (qr matrix) toString(trimRight bool) string {
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

func qr(decoder bool) []Run {
	content, qr := getStandardQr()
	qrString := qr.toString(!decoder)

	if decoder {
		return []Run{{Args: []string{qrString}, Answer: content}}
	}

	return []Run{{
		Args:   []string{content + " " + hex.EncodeToString(qr.getErrorCorrectionBlocks())},
		Answer: qrString,
	}}
}
