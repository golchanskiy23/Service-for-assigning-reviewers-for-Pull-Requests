package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

func TestSendError(t *testing.T) {
	tests := []struct {
		name           string
		code           entity.ErrorCode
		message        string
		expectedCode   entity.ErrorCode
		expectedMsg    string
		status         int
		expectedStatus int
	}{
		{
			name:           "send NOT_FOUND error",
			status:         http.StatusNotFound,
			code:           entity.CodeNotFound,
			message:        "resource not found",
			expectedStatus: http.StatusNotFound,
			expectedCode:   entity.CodeNotFound,
			expectedMsg:    "resource not found",
		},
		{
			name:           "send TEAM_EXISTS error",
			status:         http.StatusBadRequest,
			code:           entity.CodeTeamExists,
			message:        "team already exists",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   entity.CodeTeamExists,
			expectedMsg:    "team already exists",
		},
		{
			name:           "send PR_EXISTS error",
			status:         http.StatusConflict,
			code:           entity.CodePRExists,
			message:        "PR already exists",
			expectedStatus: http.StatusConflict,
			expectedCode:   entity.CodePRExists,
			expectedMsg:    "PR already exists",
		},
		{
			name:           "send PR_MERGED error",
			status:         http.StatusConflict,
			code:           entity.CodePRMerged,
			message:        "PR already merged",
			expectedStatus: http.StatusConflict,
			expectedCode:   entity.CodePRMerged,
			expectedMsg:    "PR already merged",
		},
		{
			name:           "send NOT_ASSIGNED error",
			status:         http.StatusConflict,
			code:           entity.CodeNotAssigned,
			message:        "reviewer not assigned",
			expectedStatus: http.StatusConflict,
			expectedCode:   entity.CodeNotAssigned,
			expectedMsg:    "reviewer not assigned",
		},
		{
			name:           "send NO_CANDIDATE error",
			status:         http.StatusConflict,
			code:           entity.CodeNoCandidate,
			message:        "no candidate available",
			expectedStatus: http.StatusConflict,
			expectedCode:   entity.CodeNoCandidate,
			expectedMsg:    "no candidate available",
		},
		{
			name:           "send internal server error",
			status:         http.StatusInternalServerError,
			code:           entity.CodeNotFound,
			message:        "internal server error",
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   entity.CodeNotFound,
			expectedMsg:    "internal server error",
		},
		{
			name:           "send bad request error",
			status:         http.StatusBadRequest,
			code:           entity.CodeNotFound,
			message:        "bad request",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   entity.CodeNotFound,
			expectedMsg:    "bad request",
		},
		{
			name:           "empty message",
			status:         http.StatusNotFound,
			code:           entity.CodeNotFound,
			message:        "",
			expectedStatus: http.StatusNotFound,
			expectedCode:   entity.CodeNotFound,
			expectedMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			SendError(w, tt.status, tt.code, tt.message)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var resp entity.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.Error.Code)
			assert.Equal(t, tt.expectedMsg, resp.Error.Message)
		})
	}
}
