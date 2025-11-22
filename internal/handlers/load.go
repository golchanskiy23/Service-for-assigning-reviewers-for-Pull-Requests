package handlers

import (
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func (s *Services) LoadTestHandler(w http.ResponseWriter, _ *http.Request) {
	const loadTestFreq = 100
	const loadTestDuration = 5 * time.Second
	rate := vegeta.Rate{Freq: loadTestFreq, Per: time.Second}
	duration := loadTestDuration
	go s.LoadService.RunLoadTest(rate, duration)
	if _, err := w.Write([]byte("Load test started")); err != nil {
		s.Log.Error("failed to write response", "error", err)
	}
}
