package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GetTimeStamp(now time.Time) string {
	year, month, day := now.Date()
	return fmt.Sprintf(
		"%4d-%02d-%02d %02d:%02d:%02d",
		year, month, day, now.Hour(), now.Minute(), now.Second())
}

func GetTimeStampNumber(now time.Time) string {
	year, month, day := now.Date()
	return fmt.Sprintf(
		"%4d%02d%02d%02d%02d%02d.%06d",
		year, month, day,
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond()/int(time.Microsecond))
}

func GetToken(length int) string {
	token := make([]byte, length)
	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!#$%^&*"
	for i:=0; i<length; i++  {
		token[i] = bytes[rand.Int()%len(bytes)]
	}
	return string(token)
}

func GetDate() string {
	now := time.Now()
	year, month, day := now.Date()
	return fmt.Sprintf("%4d-%02d-%02d", year, month, day)
}

func Rand(max int) int {
	return rand.Int()%max
}

func init()  {
	rand.Seed(time.Now().Unix())
}