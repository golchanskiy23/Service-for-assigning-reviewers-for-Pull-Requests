package util

import (
	"math"
	"math/rand"
	"time"
)

const (
	baseDelay    = time.Millisecond * 50
	exponentBase = 2
)

//nolint:gosec // math/rand is sufficient for jitter, crypto/rand is not needed
func CreateNewDelay(attempt int, maxVal time.Duration) time.Duration {
	backoff := baseDelay *
		time.Duration(math.Pow(exponentBase, float64(attempt)))
	if backoff > maxVal {
		backoff = maxVal
	}

	return time.Duration(rand.Int63n(int64(backoff)))
}
