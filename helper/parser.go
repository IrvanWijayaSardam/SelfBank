package helper

import "strconv"

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func StringToUint64(str string) uint64 {
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}

	return value
}

func StringToInt64(str string) int64 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}

	return value
}
