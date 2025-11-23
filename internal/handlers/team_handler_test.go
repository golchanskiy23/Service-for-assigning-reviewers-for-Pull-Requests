package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTeamService struct {
	mock.Mock
}

func (m *MockTeamService) AddTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	args := m.Called(ctx, team)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Team), args.Error(1)
}

func (m *MockTeamService) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Team), args.Error(1)
}

func TestServices_TeamAddHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*MockTeamService)
		expectedStatus int
		expectedError  bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful team creation",
			requestBody: entity.Team{
				TeamName: "team1",
				Members: []entity.TeamMember{
					{UserID: "user1", Username: "user1", IsActive: true},
					{UserID: "user2", Username: "user2", IsActive: true},
				},
			},
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("AddTeam", mock.Anything, mock.AnythingOfType("*entity.Team")).Return(&entity.Team{
					TeamName: "team1",
					Members: []entity.TeamMember{
						{UserID: "user1", Username: "user1", IsActive: true},
						{UserID: "user2", Username: "user2", IsActive: true},
					},
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp TeamAddResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "team1", resp.Team.TeamName)
				assert.Equal(t, 2, len(resp.Team.Members))
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			setupMocks:     func(teamService *MockTeamService) {},
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
			name: "empty team name",
			requestBody: entity.Team{
				TeamName: "",
				Members:  []entity.TeamMember{},
			},
			setupMocks:     func(teamService *MockTeamService) {},
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
			name: "team already exists",
			requestBody: entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("AddTeam", mock.Anything, mock.AnythingOfType("*entity.Team")).Return(nil, errors.New("TEAM_EXISTS"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, entity.CodeTeamExists, resp.Error.Code)
			},
		},
		{
			name: "add team error",
			requestBody: entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("AddTeam", mock.Anything, mock.AnythingOfType("*entity.Team")).Return(nil, errors.New("db error"))
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
			name: "team with no members",
			requestBody: entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("AddTeam", mock.Anything, mock.AnythingOfType("*entity.Team")).Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp TeamAddResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "team1", resp.Team.TeamName)
				assert.Equal(t, 0, len(resp.Team.Members))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			teamService := new(MockTeamService)
			tt.setupMocks(teamService)

			services := &Services{
				TeamService: teamService,
			}

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			services.TeamAddHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			teamService.AssertExpectations(t)
		})
	}
}

func TestServices_TeamGetHandler(t *testing.T) {
	tests := []struct {
		name           string
		teamName       string
		setupMocks     func(*MockTeamService)
		expectedStatus int
		expectedError  bool
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:     "successful get team",
			teamName: "team1",
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members: []entity.TeamMember{
						{UserID: "user1", Username: "user1", IsActive: true},
						{UserID: "user2", Username: "user2", IsActive: true},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.Team
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "team1", resp.TeamName)
				assert.Equal(t, 2, len(resp.Members))
			},
		},
		{
			name:           "missing team_name parameter",
			teamName:       "",
			setupMocks:     func(teamService *MockTeamService) {},
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
			name:     "team not found",
			teamName: "team1",
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("GetTeam", mock.Anything, "team1").Return(nil, errors.New("NOT_FOUND"))
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
			name:     "get team error",
			teamName: "team1",
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("GetTeam", mock.Anything, "team1").Return(nil, errors.New("db error"))
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
			name:     "team with no members",
			teamName: "team1",
			setupMocks: func(teamService *MockTeamService) {
				teamService.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp entity.Team
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "team1", resp.TeamName)
				assert.Equal(t, 0, len(resp.Members))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			teamService := new(MockTeamService)
			tt.setupMocks(teamService)

			services := &Services{
				TeamService: teamService,
			}

			req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+tt.teamName, nil)
			w := httptest.NewRecorder()

			services.TeamGetHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			teamService.AssertExpectations(t)
		})
	}
}
