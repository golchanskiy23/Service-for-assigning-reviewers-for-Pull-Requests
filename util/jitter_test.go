package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewDelay(t *testing.T) {
	tests := []struct {
		checkFunc func(*testing.T, time.Duration)
		name      string
		attempt   int
		maxVal    time.Duration
	}{
		{
			name:    "attempt 0, max 1s (defaultMaxConnectTimeout)",
			attempt: 0,
			maxVal:  time.Second,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*50)
			},
		},
		{
			name:    "attempt 1, max 1s (defaultMaxConnectTimeout)",
			attempt: 1,
			maxVal:  time.Second,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*200)
			},
		},
		{
			name:    "attempt 2, max 1s",
			attempt: 2,
			maxVal:  time.Second,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*400)
			},
		},
		{
			name:    "attempt 3, max 1s",
			attempt: 3,
			maxVal:  time.Second,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*800)
			},
		},
		{
			name:    "attempt 4, max 1s (capped)",
			attempt: 4,
			maxVal:  time.Second,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Second)
			},
		},
		{
			name:    "attempt 5, max 1s (capped)",
			attempt: 5,
			maxVal:  time.Second,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Second)
			},
		},
		{
			name:    "attempt 0, max 500ms",
			attempt: 0,
			maxVal:  time.Millisecond * 500,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*100)
			},
		},
		{
			name:    "attempt 1, max 500ms",
			attempt: 1,
			maxVal:  time.Millisecond * 500,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*200)
			},
		},
		{
			name:    "attempt 2, max 500ms",
			attempt: 2,
			maxVal:  time.Millisecond * 500,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*400)
			},
		},
		{
			name:    "attempt 3, max 500ms (capped)",
			attempt: 3,
			maxVal:  time.Millisecond * 500,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*500)
			},
		},
		{
			name:    "attempt 0, max 50ms",
			attempt: 0,
			maxVal:  time.Millisecond * 50,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Millisecond*50)
			},
		},
		{
			name:    "attempt 10, max 1s (capped)",
			attempt: 10,
			maxVal:  time.Second,
			checkFunc: func(t *testing.T, delay time.Duration) {
				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := CreateNewDelay(tt.attempt, tt.maxVal)
			tt.checkFunc(t, delay)
		})
	}
}

func TestCreateNewDelay_MultipleCalls(t *testing.T) {
	attempt := 2
	maxVal := time.Second
	delays := make(map[time.Duration]bool)

	for range 100 {
		delay := CreateNewDelay(attempt, maxVal)
		delays[delay] = true
	}

	assert.Greater(t, len(delays), 1, "Expected multiple different delay values due to randomness")
}

func TestCreateNewDelay_BackoffCalculation(t *testing.T) {
	maxVal := time.Second * 10

	delays := make([]time.Duration, 5)
	for i := range 5 {
		delays[i] = CreateNewDelay(i, maxVal)
	}

	assert.LessOrEqual(t, delays[0], time.Millisecond*50)
	assert.LessOrEqual(t, delays[1], time.Millisecond*100)
	assert.LessOrEqual(t, delays[2], time.Millisecond*200)
	assert.LessOrEqual(t, delays[3], time.Millisecond*400)
	assert.LessOrEqual(t, delays[4], time.Millisecond*800)
}
