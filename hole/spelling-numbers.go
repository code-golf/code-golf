package hole

import (
	"math/rand"
	"strconv"
)

var teens = [][]byte{
	[]byte("zero"),
	[]byte("one"),
	[]byte("two"),
	[]byte("three"),
	[]byte("four"),
	[]byte("five"),
	[]byte("six"),
	[]byte("seven"),
	[]byte("eight"),
	[]byte("nine"),
	[]byte("ten"),
	[]byte("eleven"),
	[]byte("twelve"),
	[]byte("thirteen"),
	[]byte("fourteen"),
	[]byte("fifteen"),
	[]byte("sixteen"),
	[]byte("seventeen"),
	[]byte("eighteen"),
	[]byte("nineteen"),
}

var tens = [][]byte{
	nil,
	nil,
	[]byte("twenty"),
	[]byte("thirty"),
	[]byte("forty"),
	[]byte("fifty"),
	[]byte("sixty"),
	[]byte("seventy"),
	[]byte("eighty"),
	[]byte("ninety"),
}

func wordify(i uint16) []byte {
	if i == 1000 {
		return []byte("one thousand")
	} else if i < 20 {
		return teens[i]
	} else if i < 100 {
		if j := i % 10; j > 0 {
			return append(append(append([]byte(nil), tens[i/10]...), '-'), teens[j]...)
		}

		return tens[i/10]
	}

	hundred := append(append([]byte(nil), teens[i/100]...), []byte(" hundred")...)

	if j := i % 100; j > 0 {
		hundred = append(append(hundred, []byte(" and ")...), wordify(j)...)
	}

	return hundred
}

func spellingNumbers() (args []string, out string) {
	// Always test the low and high end.
	numbers := []uint16{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		990, 991, 992, 993, 994, 995, 996, 997, 998, 999, 1000,
	}

	tens := []uint16{20, 30, 40, 50, 60, 70, 80, 90}

	rand.Shuffle(len(tens), func(i, j int) {
		tens[i], tens[j] = tens[j], tens[i]
	})

	// Add on random numbers to each ten except the first, to ensure at least
	// one non hyphen answer.
	for i := 1; i < len(tens); i++ {
		tens[i] += 1 + uint16(rand.Intn(8))
	}

	numbers = append(numbers, tens...)

	hundreds := []uint16{100, 200, 300, 400, 500, 600, 700, 800, 900}

	rand.Shuffle(len(hundreds), func(i, j int) {
		hundreds[i], hundreds[j] = hundreds[j], hundreds[i]
	})

	// Add on random numbers to each hundred except the first, to ensure at
	// least one non and answer.
	for i := 1; i < len(hundreds); i++ {
		hundreds[i] += 1 + uint16(rand.Intn(98))
	}

	numbers = append(numbers, hundreds...)

	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	for _, n := range numbers {
		args = append(args, strconv.Itoa(int(n)))
		out += string(wordify(n)) + "\n"
	}

	// Drop the trailing newline.
	out = out[:len(out)-1]

	return
}
