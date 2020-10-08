package hole

import (
	"math/rand"
	"strconv"
	"strings"
)

// a bounding box (bbox) is defined in
// terms of its top-left vertex coordinates
// (x, y) and its width and height (w, h).
type bbox struct{ x, y, w, h int }

// couldn't find a quick way to loop a struct
func strconvbox(box bbox) (out string) {
	var outs []string
	outs = append(outs, strconv.Itoa(box.x))
	outs = append(outs, strconv.Itoa(box.y))
	outs = append(outs, strconv.Itoa(box.w))
	outs = append(outs, strconv.Itoa(box.h))
	return strings.Join(outs, " ")
}

// compute bottom-right (br) bbox coordinates
func unbox(b bbox) (tlx, tly, brx, bry int) {
	tlx = b.x
	tly = b.y
	brx = b.x + b.w
	bry = b.y + b.h
	return
}

// til that go doesn't have built-in max/min for int
func minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func calculateIntersection(b1, b2 bbox) int {
	tlx1, tly1, brx1, bry1 := unbox(b1)
	tlx2, tly2, brx2, bry2 := unbox(b2)

	// find top-left and bottom-right intersection coordinates
	itlx := maxint(tlx1, tlx2) // intersection top left x
	itly := maxint(tly1, tly2)
	ibrx := minint(brx1, brx2) // intersection bottom right x
	ibry := minint(bry1, bry2)

	// calculate intersection dimensions
	iw := ibrx - itlx
	ih := ibry - itly

	// intersection is empty if the bboxes do not overlap
	if iw < 0 || ih < 0 || iw > b1.w+b2.w || ih > b1.h+b2.h {
		return 0
	}
	return iw * ih
}

// generator of random non-null boxes (i.e. with area != 0)
func boxGen() bbox {
	return bbox{
		x: rand.Intn(101),
		y: rand.Intn(101),
		w: rand.Intn(50) + 1,
		h: rand.Intn(50) + 1,
	}
}

func intersection() (args []string, out string) {
	var outs []string

	//// default cases
	// define two non overlapping 1x1 boxes
	b1 := bbox{x: 0, y: 0, h: 1, w: 1}
	b2 := bbox{x: 0, y: 0, h: 2, w: 2}
	b3 := bbox{x: 3, y: 3, h: 1, w: 2}
	b4 := bbox{x: 3, y: 1, h: 3, w: 2}
	b5 := bbox{x: 3, y: 1, h: 3, w: 1}
	b6 := bbox{x: 0, y: 0, h: 10, w: 10}
	b7 := bbox{x: 2, y: 2, h: 2, w: 2}

	// b1 and b2 overlap by 1 pixel
	args = append(args, strconvbox(b1)+" "+strconvbox(b2))
	outs = append(outs, "1")

	// b1 and b3 are far away and don't overlap
	args = append(args, strconvbox(b1)+" "+strconvbox(b3))
	outs = append(outs, "0")

	// b3 and b4 overlap on one horizontal side
	args = append(args, strconvbox(b3)+" "+strconvbox(b4))
	outs = append(outs, "2")

	// b4 and b5 overlap on one vertical side
	args = append(args, strconvbox(b4)+" "+strconvbox(b5))
	outs = append(outs, "3")

	// b4 is inside b6
	args = append(args, strconvbox(b4)+" "+strconvbox(b6))
	outs = append(outs, "6")

	// b2 and b7 are side by side but don't overlap
	args = append(args, strconvbox(b2)+" "+strconvbox(b7))
	outs = append(outs, "0")

	//// generate 100 random cases
	zeros := 0
	nonZeros := 0
	for zeros+nonZeros < 100 {
		b1 = boxGen()
		b2 = boxGen()
		intersection := calculateIntersection(b1, b2)

		// compute 90 non-zero cases and 10 zero ones
		if intersection > 0 && nonZeros < 90 {
			args = append(args, strconvbox(b1)+" "+strconvbox(b2))
			outs = append(outs, strconv.Itoa(intersection))
			nonZeros++
		} else if intersection == 0 && zeros < 10 {
			args = append(args, strconvbox(b1)+" "+strconvbox(b2))
			outs = append(outs, strconv.Itoa(intersection))
			zeros++
		}
	}

	// 13x13 default side cases
	// - |   |
	// --|   |
	// --|-  |
	// --|---|
	// --|---|-
	//   |-  |
	//   |---|
	//   |---|-
	//   | - |
	//   | --|
	//   | --|-
	//   |   |-
	//   |   | -
	bigbox := bbox{x: 2, y: 2, w: 3, h: 3}
	strbigbox := strconvbox(bigbox)

	xs := []int{0, 2, 3, 5, 6}
	ys := []int{0, 2, 3, 5, 6}
	for _, x := range xs {
		for _, y := range ys {
			for w := 1; w < 7; w++ {
				for h := 1; h < 7; h++ {
					if (x == 0 && (w == 4 || w > 6)) ||
						(x == 2 && (w == 2 || w > 4)) ||
						(x == 3 && w > 3) ||
						(x == 5 && w > 1) || (x == 6 && w > 1) {
						continue
					}
					if (y == 0 && (h == 4 || h > 6)) ||
						(y == 2 && (h == 2 || h > 4)) ||
						(y == 3 && h > 3) ||
						(y == 5 && h > 1) || (y == 6 && h > 1) {
						continue
					}
					if rand.Float32() > 0.5 { // randomly add test ?
						b := bbox{x: x, y: y, w: w, h: h}
						if rand.Float32() > 0.5 { // randomly flip input
							args = append(args, strconvbox(b)+" "+strbigbox)
						} else {
							args = append(args, strbigbox+" "+strconvbox(b))
						}
						outs = append(outs, strconv.Itoa(calculateIntersection(b, bigbox)))
					}
				}
			}
		}
	}

	// shuffle args and outputs in the same way
	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	out = strings.Join(outs, "\n")
	return
}
