package handlers

import (
	"context"
	"net/http"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	assignedReviewersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "assigned_reviewers_per_pr",
			Help: "Number of assigned reviewers for a given open pull request",
		},
		[]string{"pull_request_id"},
	)

	openPRsPerUserGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "open_prs_per_user",
			Help: "Number of open pull requests currently assigned to a user",
		},
		[]string{"user_id"},
	)
)

var metricsRegistry = prometheus.NewRegistry()

func init() {
	metricsRegistry.MustRegister(assignedReviewersGauge, openPRsPerUserGauge)
}

func (s *Services) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	prCounts, err := s.StatsService.GetAssignedCountPerPR(ctx)
	if err != nil {
		http.Error(w, "failed to get PR counts", http.StatusInternalServerError)
		return
	}

	userCounts, err := s.StatsService.GetOpenPRCountPerUser(ctx)
	if err != nil {
		util.SendError(w,
			http.StatusInternalServerError,
			entity.CodePRCount,
			"failed to get user counts")
		return
	}

	assignedReviewersGauge.Reset()
	openPRsPerUserGauge.Reset()

	for prID, cnt := range prCounts {
		assignedReviewersGauge.WithLabelValues(prID).Set(float64(cnt))
	}

	for userID, cnt := range userCounts {
		openPRsPerUserGauge.WithLabelValues(userID).Set(float64(cnt))
	}

	promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func (s *Services) UpdateMetrics(ctx context.Context) error {
	prCounts, err := s.StatsService.GetAssignedCountPerPR(ctx)
	if err != nil {
		return err
	}

	userCounts, err := s.StatsService.GetOpenPRCountPerUser(ctx)
	if err != nil {
		return err
	}

	assignedReviewersGauge.Reset()
	openPRsPerUserGauge.Reset()

	for prID, cnt := range prCounts {
		assignedReviewersGauge.WithLabelValues(prID).Set(float64(cnt))
	}

	for userID, cnt := range userCounts {
		openPRsPerUserGauge.WithLabelValues(userID).Set(float64(cnt))
	}

	return nil
}
