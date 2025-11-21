package util

import (
	"math"
	"math/rand"
	"time"
)

const baseDelay = time.Millisecond * 50

func CreateNewDelay(attempt int, maxVal time.Duration) time.Duration {
	backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
	if backoff > maxVal {
		backoff = maxVal
	}
	return time.Duration(rand.Int63n(int64(backoff)))
}
