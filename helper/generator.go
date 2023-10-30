package helper

import (
	"math/rand"
	"time"
)

func GenerateRandomAccountNumber() uint64 {
	rand.Seed(time.Now().UnixNano())
	min := 1000000 // 7-digit minimum number
	max := 9999999 // 7-digit maximum number
	return uint64(min + rand.Intn(max-min+1))
}

func GenerateTrxId() uint64 {
	rand.Seed(time.Now().UnixNano())
	min := 1000000000 // 7-digit minimum number
	max := 9999999999 // 7-digit maximum number
	return uint64(min + rand.Intn(max-min+1))
}

func GetCurrentTimeInLocation() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}
	}

	currentTime := time.Now().In(loc)

	return currentTime
}
