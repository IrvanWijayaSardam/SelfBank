package helper

import "strconv"

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}
