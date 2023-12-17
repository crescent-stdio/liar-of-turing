package utils

import (
	"math/rand"
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetNewRand() *rand.Rand {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rnd := rand.New(src)
	return rnd
}

func RandomTimeSleep() {
	rnd := GetNewRand()
	duration := rnd.Intn(1000) + 800
	time.Sleep(time.Duration(duration) * time.Millisecond)
}
