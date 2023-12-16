package utils

import (
	"math/rand"
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetNewRand() *rand.Rand {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rnd := rand.New(src)
	return rnd
}
