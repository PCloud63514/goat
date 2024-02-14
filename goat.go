package goat

import "time"

var (
	startUpDateTime = time.Now()
)

func GetStartUpDateTime() time.Time {
	return startUpDateTime
}
