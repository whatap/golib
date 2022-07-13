package hexa32

import (
	"math"
	"strconv"
)

const PLUS = 'x'
const MINUS = 'z'

func ToString32(num int64) string {
	minus := num < 0
	if minus {
		if num == math.MinInt64 {
			return "z8000000000000"
		}
		return "z" + to_str(-num)
	} else {
		if num < 10 {
			return strconv.Itoa(int(num))
		} else {
			return "x" + to_str(num)
		}
	}
}
func ToLong32(str string) int64 {
	if str == "" {
		return 0
	}
	switch str[0] {
	case MINUS:
		if "z8000000000000" == str {
			return math.MinInt64
		} else {
			return -1 * to_long(str[1:len(str)])
		}
	case PLUS:
		return to_long(str[1:len(str)])
	default:
		i, err := strconv.Atoi(str)
		if err != nil {
			return 0
		}
		return int64(i)
	}

}

var digits = []byte{
	'0', '1', '2', '3', '4', '5',
	'6', '7', '8', '9', 'a', 'b',
	'c', 'd', 'e', 'f', 'g', 'h',
	'i', 'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z'}

func to_str(i int64) string {
	radix := int64(32)
	buf := make([]byte, 65)
	charPos := 64
	for i = -i; i <= (-radix); i = i / radix {
		buf[charPos] = digits[int(-(i % radix))]
		charPos--
	}
	buf[charPos] = digits[int(-i)]

	return string(buf[charPos:65])
}

func to_long(s string) int64 {
	var result int64 = 0
	var limit int64 = -math.MaxInt64
	var multmin int64 = limit / 32
	sz := len(s)

	findc := func(x int) int64 {
		switch {
		case '0' <= x && x <= '9':
			return int64(x - '0')
		case 'a' <= x && x <= 'z':
			return int64(x - 'a' + 10)
		case 'A' <= x && x <= 'Z':
			return int64(x - 'A' + 10)
		default:
			return 0
		}
	}

	for i := 0; i < sz; i++ {
		digit := findc(int(s[i]))
		if result < multmin {
			return 0
		}
		result *= 32
		if result < limit+digit {
			return 0
		}
		result -= digit
	}
	return -result
}
