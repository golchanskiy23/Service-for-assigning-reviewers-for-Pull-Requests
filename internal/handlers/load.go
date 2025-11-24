package handlers

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"net/http"
	"time"
)

func (s *Services) LoadTestHandler(w http.ResponseWriter, r *http.Request) {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 5 * time.Second
	go s.LoadService.RunLoadTest(rate, duration)
	w.Write([]byte("Load test started"))
}
