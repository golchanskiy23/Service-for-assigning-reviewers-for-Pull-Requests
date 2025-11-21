package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/handlers"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) ChangeStatus(
	ctx context.Context,
	userID string,
	isActive bool,
) (*entity.User, error) {
	args := m.Called(ctx, userID, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	user, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, args.Error(1)
	}

	return user, args.Error(1)
}

func (m *MockUserService) GetPRsAssignedTo(
	ctx context.Context,
	userID string,
) (string, []*entity.PullRequestShort, error) {
	args := m.Called(ctx, userID)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}

	prs, ok := args.Get(1).([]*entity.PullRequestShort)
	if !ok {
		return args.String(0), nil, args.Error(2)
	}

	return args.String(0), prs, args.Error(2)
}

func (m *MockUserService) MassDeactivate(ctx context.Context, users []entity.User, flag bool) error {
	args := m.Called(ctx, users, flag)
	return args.Error(0)
}

func TestServices_UserSetIsActiveHandler(t *testing.T) {
	tests := []struct {
		requestBody    interface{}
		setupMocks     func(*MockUserService)
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
		name           string
		expectedStatus int
		expectedError  bool
	}{
		//nolint:dupl // necessary tests
		{
			name: "successful activation",
			requestBody: handlers.UserSetIsActiveRequest{
				UserID:   "user1",
				IsActive: true,
			},
			setupMocks: func(userService *MockUserService) {
				userService.On("ChangeStatus", mock.Anything, "user1", true).Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp handlers.UserSetIsActiveResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "user1", resp.User.UserID)
				assert.True(t, resp.User.IsActive)
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			setupMocks:     func(userService *MockUserService) {},
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
			name: "user not found",
			requestBody: handlers.UserSetIsActiveRequest{
				UserID:   "user1",
				IsActive: true,
			},
			setupMocks: func(userService *MockUserService) {
				userService.On("ChangeStatus", mock.Anything, "user1", true).Return(nil, errors.New("NOT_FOUND"))
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
			name: "change status error",
			requestBody: handlers.UserSetIsActiveRequest{
				UserID:   "user1",
				IsActive: true,
			},
			setupMocks: func(userService *MockUserService) {
				userService.On("ChangeStatus", mock.Anything, "user1", true).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
		//nolint:dupl // necessary tests
		{
			name: "successful deactivation",
			requestBody: handlers.UserSetIsActiveRequest{
				UserID:   "user1",
				IsActive: false,
			},
			setupMocks: func(userService *MockUserService) {
				userService.On("ChangeStatus", mock.Anything, "user1", false).Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp handlers.UserSetIsActiveResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "user1", resp.User.UserID)
				assert.False(t, resp.User.IsActive)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := new(MockUserService)
			tt.setupMocks(userService)

			services := &handlers.Services{
				UserService: userService,
			}

			var body []byte

			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/user/set-active", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			services.UserSetIsActiveHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			userService.AssertExpectations(t)
		})
	}
}

func TestServices_UserGetReviewHandler(t *testing.T) {
	tests := []struct {
		setupMocks     func(*MockUserService)
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
		name           string
		userID         string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful get PRs",
			userID: "user1",
			setupMocks: func(userService *MockUserService) {
				userService.On("GetPRsAssignedTo", mock.Anything, "user1").Return("user1", []*entity.PullRequestShort{
					{
						PullRequestID:   "pr1",
						PullRequestName: "Test PR 1",
						AuthorID:        "user2",
						Status:          entity.OPEN,
					},
					{
						PullRequestID:   "pr2",
						PullRequestName: "Test PR 2",
						AuthorID:        "user3",
						Status:          entity.MERGED,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp handlers.UserGetReviewResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "user1", resp.UserID)
				assert.Len(t, resp.PullRequests, 2)
			},
		},
		{
			name:           "missing user_id parameter",
			userID:         "",
			setupMocks:     func(userService *MockUserService) {},
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
			name:   "user not found",
			userID: "user1",
			setupMocks: func(userService *MockUserService) {
				userService.On("GetPRsAssignedTo", mock.Anything, "user1").Return("", nil, errors.New("NOT_FOUND"))
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
			name:   "get PRs error",
			userID: "user1",
			setupMocks: func(userService *MockUserService) {
				userService.On("GetPRsAssignedTo", mock.Anything, "user1").Return("", nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
			},
		},
		{
			name:   "no PRs assigned",
			userID: "user1",
			setupMocks: func(userService *MockUserService) {
				userService.On("GetPRsAssignedTo", mock.Anything, "user1").Return("user1", []*entity.PullRequestShort{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp handlers.UserGetReviewResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "user1", resp.UserID)
				assert.Empty(t, resp.PullRequests)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := new(MockUserService)
			tt.setupMocks(userService)

			services := &handlers.Services{
				UserService: userService,
			}

			req := httptest.NewRequest(http.MethodGet, "/user/review?user_id="+tt.userID, http.NoBody)
			w := httptest.NewRecorder()

			services.UserGetReviewHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			userService.AssertExpectations(t)
		})
	}
}
