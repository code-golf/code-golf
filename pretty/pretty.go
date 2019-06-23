package pretty

import (
	"fmt"
	"strconv"
)

// NOTE Only handles 0 - 999,999
func Comma(i int) string {
	if i > 999 {
		return fmt.Sprintf("%d,%03d", i/1000, i%1000)
	}

	return strconv.Itoa(i)
}
