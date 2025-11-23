package integration

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/server"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLoadService struct {
	mock.Mock
}

func (m *MockLoadService) RunLoadTest(rate vegeta.Rate, duration time.Duration) {
	m.Called(rate, duration)
}

func setupRouterWithServices(s *handlers.Services) *chi.Mux {
	r := chi.NewRouter()
	server.RegisterRoutes(s, r)
	return r
}

func TestLoadTestHandler(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		setupMock       func(*MockLoadService)
		expectedCode    int
		expectedBodySub string
	}{
		{
			name:  "valid request",
			query: "?freq=10&duration=3s",
			setupMock: func(m *MockLoadService) {
				m.On("RunLoadTest", mock.Anything, mock.Anything).Return()
			},
			expectedCode:    http.StatusOK,
			expectedBodySub: "Load test started",
		},
		{
			name:            "missing freq",
			query:           "?duration=3s",
			setupMock:       func(m *MockLoadService) {},
			expectedCode:    http.StatusBadRequest,
			expectedBodySub: "freq is required",
		},
		{
			name:            "invalid duration",
			query:           "?freq=10&duration=abc",
			setupMock:       func(m *MockLoadService) {},
			expectedCode:    http.StatusBadRequest,
			expectedBodySub: "duration must be valid duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoad := &MockLoadService{}
			if tt.setupMock != nil {
				tt.setupMock(mockLoad)
			}

			s := &handlers.Services{
				LoadService: mockLoad,
				Log:         nil,
			}

			if tt.name == "LoadService nil" {
				s.LoadService = nil
			}

			router := setupRouterWithServices(s)
			req := httptest.NewRequest(http.MethodGet, "/loadtest"+tt.query, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBodySub)

		})
	}
}
