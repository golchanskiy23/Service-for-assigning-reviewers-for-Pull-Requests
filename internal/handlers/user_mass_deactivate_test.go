package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

func TestServices_UserMassDeactivateHandler(t *testing.T) {
	t.Run("successful deactivation", func(t *testing.T) {
		userService := new(MockUserService)
		userService.On("MassDeactivate",
			mock.Anything,
			[]entity.User{{UserID: "u1",
				Username: "Alice",
				TeamName: "backend",
				IsActive: true},
				{UserID: "u2",
					Username: "Bob",
					TeamName: "backend",
					IsActive: false}}, false).Return(nil)

		services := &Services{UserService: userService}

		reqBody := map[string]interface{}{
			"users": []map[string]interface{}{
				{"user_id": "u1",
					"username":  "Alice",
					"team_name": "backend",
					"is_active": true},
				{"user_id": "u2",
					"username":  "Bob",
					"team_name": "backend",
					"is_active": false},
			},
			"flag": false,
		}

		b, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users/deactivate", bytes.NewBuffer(b))
		w := httptest.NewRecorder()

		services.UsersMassDeactivateHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp UserMassChangeResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"u1", "u2"}, resp.Deactivated)

		userService.AssertExpectations(t)
	})

	t.Run("invalid json returns bad request", func(t *testing.T) {
		userService := new(MockUserService)
		services := &Services{UserService: userService}

		req := httptest.NewRequest(http.MethodPost, "/users/deactivate", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()

		services.UsersMassDeactivateHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("flag true returns bad request", func(t *testing.T) {
		userService := new(MockUserService)
		services := &Services{UserService: userService}

		reqBody := map[string]interface{}{
			"users": []map[string]interface{}{
				{"user_id": "u1", "team_name": "backend"}}, "flag": true}
		b, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users/deactivate", bytes.NewBuffer(b))
		w := httptest.NewRecorder()

		services.UsersMassDeactivateHandler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns NOT_FOUND -> 404", func(t *testing.T) {
		userService := new(MockUserService)
		userService.On("MassDeactivate", mock.Anything, mock.Anything, false).Return(errors.New("NOT_FOUND"))

		services := &Services{UserService: userService}

		reqBody := map[string]interface{}{
			"users": []map[string]interface{}{
				{"user_id": "u1", "team_name": "backend"}}, "flag": false}
		b, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users/deactivate", bytes.NewBuffer(b))
		w := httptest.NewRecorder()

		services.UsersMassDeactivateHandler(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		userService.AssertExpectations(t)
	})
}
