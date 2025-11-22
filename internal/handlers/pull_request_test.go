package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockPRService struct {
	mock.Mock
}

func (m *MockPRService) CreatePR(ctx context.Context, prID, prName, authorID string) (*entity.PullRequest, string, error) {
	args := m.Called(ctx, prID, prName, authorID)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*entity.PullRequest), args.String(1), args.Error(2)
}

func (m *MockPRService) MergePR(ctx context.Context, prID string) (*entity.PullRequest, error) {
	args := m.Called(ctx, prID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.PullRequest), args.Error(1)
}

func (m *MockPRService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*entity.PullRequest, string, error) {
	args := m.Called(ctx, prID, oldReviewerID)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*entity.PullRequest), args.String(1), args.Error(2)
}

func TestServices_PRCreateHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*MockPRService)
		expectedStatus int
		expectedError  bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful PR creation",
			requestBody: PRCreateRequest{
				PullRequestID:   "pr1",
				PullRequestName: "Test PR",
				AuthorID:        "user1",
			},
			setupMocks: func(prService *MockPRService) {
				now := time.Now()
				prService.On("CreatePR", mock.Anything, "pr1", "Test PR", "user1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
				}, "", nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp PRCreateResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "pr1", resp.PR.PullRequestID)
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			setupMocks:     func(prService *MockPRService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodeNotFound, resp.Error.Code)
			},
		},
		{
			name: "PR_EXISTS error",
			requestBody: PRCreateRequest{
				PullRequestID:   "pr1",
				PullRequestName: "Test PR",
				AuthorID:        "user1",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("CreatePR", mock.Anything, "pr1", "Test PR", "user1").Return(nil, "", errors.New("PR_EXISTS"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodePRExists, resp.Error.Code)
			},
		},
		{
			name: "NOT_FOUND error",
			requestBody: PRCreateRequest{
				PullRequestID:   "pr1",
				PullRequestName: "Test PR",
				AuthorID:        "user1",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("CreatePR", mock.Anything, "pr1", "Test PR", "user1").Return(nil, "", errors.New("NOT_FOUND"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodeNotFound, resp.Error.Code)
			},
		},
		{
			name: "internal server error",
			requestBody: PRCreateRequest{
				PullRequestID:   "pr1",
				PullRequestName: "Test PR",
				AuthorID:        "user1",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("CreatePR", mock.Anything, "pr1", "Test PR", "user1").Return(nil, "", errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prService := new(MockPRService)
			tt.setupMocks(prService)

			services := &Services{
				PRService: prService,
			}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/pr/create", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			services.PRCreateHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			prService.AssertExpectations(t)
		})
	}
}

func TestServices_PRMergeHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*MockPRService)
		expectedStatus int
		expectedError  bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful PR merge",
			requestBody: PRMergeRequest{
				PullRequestID: "pr1",
			},
			setupMocks: func(prService *MockPRService) {
				now := time.Now()
				mergedTime := time.Now()
				prService.On("MergePR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.MERGED,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
					MergedAt:          &mergedTime,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp PRMergeResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.MERGED, resp.PR.Status)
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			setupMocks:     func(prService *MockPRService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
		{
			name: "PR not found",
			requestBody: PRMergeRequest{
				PullRequestID: "pr1",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("MergePR", mock.Anything, "pr1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodeNotFound, resp.Error.Code)
			},
		},
		{
			name: "merge error",
			requestBody: PRMergeRequest{
				PullRequestID: "pr1",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("MergePR", mock.Anything, "pr1").Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prService := new(MockPRService)
			tt.setupMocks(prService)

			services := &Services{
				PRService: prService,
			}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/pr/merge", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			services.PRMergeHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			prService.AssertExpectations(t)
		})
	}
}

func TestServices_PRReassignHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*MockPRService)
		expectedStatus int
		expectedError  bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful reassignment",
			requestBody: PRReassignRequest{
				PullRequestID: "pr1",
				OldUserID:     "user2",
			},
			setupMocks: func(prService *MockPRService) {
				now := time.Now()
				prService.On("ReassignReviewer", mock.Anything, "pr1", "user2").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user3"},
					CreatedAt:         &now,
				}, "user3", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp PRReassignResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "user3", resp.ReplacedBy)
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			setupMocks:     func(prService *MockPRService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
		{
			name: "PR or user not found",
			requestBody: PRReassignRequest{
				PullRequestID: "pr1",
				OldUserID:     "user2",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("ReassignReviewer", mock.Anything, "pr1", "user2").Return(nil, "", errors.New("NOT_FOUND"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodeNotFound, resp.Error.Code)
			},
		},
		{
			name: "PR already merged",
			requestBody: PRReassignRequest{
				PullRequestID: "pr1",
				OldUserID:     "user2",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("ReassignReviewer", mock.Anything, "pr1", "user2").Return(nil, "", errors.New("PR_MERGED"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodePRMerged, resp.Error.Code)
			},
		},
		{
			name: "reviewer not assigned",
			requestBody: PRReassignRequest{
				PullRequestID: "pr1",
				OldUserID:     "user2",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("ReassignReviewer", mock.Anything, "pr1", "user2").Return(nil, "", errors.New("NOT_ASSIGNED"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodeNotAssigned, resp.Error.Code)
			},
		},
		{
			name: "no candidate available",
			requestBody: PRReassignRequest{
				PullRequestID: "pr1",
				OldUserID:     "user2",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("ReassignReviewer", mock.Anything, "pr1", "user2").Return(nil, "", errors.New("NO_CANDIDATE"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodeNoCandidate, resp.Error.Code)
			},
		},
		{
			name: "reassignment error",
			requestBody: PRReassignRequest{
				PullRequestID: "pr1",
				OldUserID:     "user2",
			},
			setupMocks: func(prService *MockPRService) {
				prService.On("ReassignReviewer", mock.Anything, "pr1", "user2").Return(nil, "", errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prService := new(MockPRService)
			tt.setupMocks(prService)

			services := &Services{
				PRService: prService,
			}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/pr/reassign", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			services.PRReassignHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			prService.AssertExpectations(t)
		})
	}
}
