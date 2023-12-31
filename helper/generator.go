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

func GetCurrentTimeInLocation() int64 {
	currentTime := time.Now()
	currentTimestamp := currentTime.Unix()
	return currentTimestamp
}

func ConvertUnixtime(epoch int64) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	date := time.Unix(epoch, 0).In(loc)

	return date
}
