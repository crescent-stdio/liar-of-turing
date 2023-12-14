package utils

import "time"

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
