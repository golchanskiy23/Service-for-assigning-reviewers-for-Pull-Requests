package integration

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/handlers"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/service"
)

type MockStatsRepo struct {
	mock.Mock
}

func (m *MockStatsRepo) GetAssignedReviewersCountPerPR(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockStatsRepo) GetOpenPRCountPerUser(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func TestMetricsHandler_ServesMetrics(t *testing.T) {
	mockRepo := new(MockStatsRepo)
	svc := service.NewStatsService(mockRepo)

	prCounts := map[string]int{"pr-1": 1, "pr-2": 2}
	userCounts := map[string]int{"u1": 1}

	mockRepo.On("GetAssignedReviewersCountPerPR", mock.Anything).Return(prCounts, nil)
	mockRepo.On("GetOpenPRCountPerUser", mock.Anything).Return(userCounts, nil)

	services := &handlers.Services{
		Log:          newTestLogger(),
		StatsService: svc}

	req := httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody)
	w := httptest.NewRecorder()

	services.MetricsHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := w.Body.String()

	assert.True(t, strings.Contains(body, "assigned_reviewers_per_pr"))
	assert.True(t, strings.Contains(body, "open_prs_per_user"))
	assert.True(t, strings.Contains(body, "pr-1"))
	assert.True(t, strings.Contains(body, "pr-2"))
	assert.True(t, strings.Contains(body, "u1"))

	mockRepo.AssertExpectations(t)
}

func TestMetricsHandler_ErrorPaths(t *testing.T) {
	t.Run("pr counts error", func(t *testing.T) {
		mockRepo := new(MockStatsRepo)
		svc := service.NewStatsService(mockRepo)

		mockRepo.On("GetAssignedReviewersCountPerPR", mock.Anything).Return(nil, errors.New("db error"))

		services := &handlers.Services{
			Log:          newTestLogger(),
			StatsService: svc}

		req := httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody)
		w := httptest.NewRecorder()

		services.MetricsHandler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "failed to get PR counts")

		mockRepo.AssertExpectations(t)
	})

	t.Run("user counts error", func(t *testing.T) {
		mockRepo := new(MockStatsRepo)
		svc := service.NewStatsService(mockRepo)

		mockRepo.On("GetAssignedReviewersCountPerPR", mock.Anything).Return(map[string]int{}, nil)
		mockRepo.On("GetOpenPRCountPerUser", mock.Anything).Return(nil, errors.New("db error"))

		services := &handlers.Services{
			Log:          newTestLogger(),
			StatsService: svc}

		req := httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody)
		w := httptest.NewRecorder()

		services.MetricsHandler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "failed to get user counts")

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateMetrics_Method(t *testing.T) {
	mockRepo := new(MockStatsRepo)
	svc := service.NewStatsService(mockRepo)

	prCounts := map[string]int{"pr-1": 1}
	userCounts := map[string]int{"u1": 1}

	mockRepo.On("GetAssignedReviewersCountPerPR", mock.Anything).Return(prCounts, nil)
	mockRepo.On("GetOpenPRCountPerUser", mock.Anything).Return(userCounts, nil)

	services := &handlers.Services{
		Log:          newTestLogger(),
		StatsService: svc}

	err := services.UpdateMetrics(context.Background())
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
