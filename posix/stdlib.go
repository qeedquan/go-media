package posix

import (
	"strconv"
	"unicode"
)

func Strtol(str string, base int) (val int64, advance int, err error) {
	for i, r := range str {
		if !(r == '+' || r == '-' || unicode.IsDigit(r)) {
			advance = i
			break
		}
	}
	val, err = strconv.ParseInt(str[:advance], base, 64)
	if err != nil {
		advance = 0
	}
	return
}
