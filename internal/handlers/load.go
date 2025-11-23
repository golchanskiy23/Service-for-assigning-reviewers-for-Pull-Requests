package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/util"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	five = 5
)

type LoadTestRequest struct {
	Freq     int           `json:"freq"`
	Duration time.Duration `json:"duration"`
}

func validateLoadTestRequest(req *LoadTestRequest) error {
	if req.Freq <= Zero {
		return errors.New("freq must be a positive integer")
	}
	if req.Duration <= Zero {
		return errors.New("duration must be a positive duration")
	}

	if req.Duration > five*time.Minute {
		return errors.New("duration is too large")
	}
	if req.Freq > 10000 {
		return errors.New("freq is too large")
	}
	return nil
}

func (s *Services) LoadTestHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println("Panic in LoadTestHandler:", rec)
			debug.PrintStack()
		}
	}()

	fmt.Println(">>> LoadTestHandler start: r nil?", r == nil)
	q := r.URL.Query()
	req := LoadTestRequest{}
	freqStr := q.Get("freq")
	if freqStr == "" {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			"freq is required",
		)

		return
	}
	freq, err := strconv.Atoi(freqStr)
	if err != nil {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			"freq must be integer",
		)

		return
	}

	req.Freq = freq
	durStr := q.Get("duration")
	if durStr == "" {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			"duration is required",
		)

		return
	}
	duration, err := time.ParseDuration(durStr)
	if err != nil {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			"duration must be valid duration like 5s or 1m",
		)
		return
	}
	req.Duration = duration

	if err := validateLoadTestRequest(&req); err != nil {
		util.SendError(w,
			http.StatusBadRequest,
			entity.CodeBadRequest,
			err.Error(),
		)
		return
	}

	rate := vegeta.Rate{Freq: req.Freq, Per: time.Second}
	go s.LoadService.RunLoadTest(rate, req.Duration)
	if _, err := w.Write([]byte("Load test started")); err != nil {
		s.Log.Error("failed to write response", "error", err)
	}
}
