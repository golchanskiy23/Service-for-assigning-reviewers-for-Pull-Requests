package service

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTeamService_AddTeam(t *testing.T) {
	tests := []struct {
		name          string
		team          *entity.Team
		setupMocks    func(*MockTeamRepository)
		expectedError string
		expectedTeam  *entity.Team
	}{
		{
			name: "successful team creation",
			team: &entity.Team{
				TeamName: "team1",
				Members: []entity.TeamMember{
					{UserID: "user1", Username: "user1", IsActive: true},
					{UserID: "user2", Username: "user2", IsActive: true},
				},
			},
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("TeamExists", mock.Anything, "team1").Return(false, nil)
				teamRepo.On("AddTeam", mock.Anything, mock.Anything).Return(nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members: []entity.TeamMember{
						{UserID: "user1", Username: "user1", IsActive: true},
						{UserID: "user2", Username: "user2", IsActive: true},
					},
				}, nil)
			},
			expectedError: "",
			expectedTeam: &entity.Team{
				TeamName: "team1",
				Members: []entity.TeamMember{
					{UserID: "user1", Username: "user1", IsActive: true},
					{UserID: "user2", Username: "user2", IsActive: true},
				},
			},
		},
		{
			name: "team already exists",
			team: &entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("TeamExists", mock.Anything, "team1").Return(true, nil)
			},
			expectedError: "TEAM_EXISTS",
			expectedTeam:  nil,
		},
		{
			name: "team exists check error",
			team: &entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("TeamExists", mock.Anything, "team1").Return(false, errors.New("db error"))
			},
			expectedError: "db error",
			expectedTeam:  nil,
		},
		{
			name: "add team error",
			team: &entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("TeamExists", mock.Anything, "team1").Return(false, nil)
				teamRepo.On("AddTeam", mock.Anything, mock.Anything).Return(errors.New("TEAM_EXISTS"))
			},
			expectedError: "TEAM_EXISTS",
			expectedTeam:  nil,
		},
		{
			name: "get team after add error",
			team: &entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("TeamExists", mock.Anything, "team1").Return(false, nil)
				teamRepo.On("AddTeam", mock.Anything, mock.Anything).Return(nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(nil, errors.New("db error"))
			},
			expectedError: "db error",
			expectedTeam:  nil,
		},
		{
			name: "team with no members",
			team: &entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("TeamExists", mock.Anything, "team1").Return(false, nil)
				teamRepo.On("AddTeam", mock.Anything, mock.Anything).Return(nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
			},
			expectedError: "",
			expectedTeam: &entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			teamRepo := new(MockTeamRepository)

			tt.setupMocks(teamRepo)

			svc := NewTeamService(teamRepo)
			ctx := context.Background()

			team, err := svc.AddTeam(ctx, tt.team)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, team)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, team)
				if tt.expectedTeam != nil {
					assert.Equal(t, tt.expectedTeam.TeamName, team.TeamName)
					assert.Equal(t, len(tt.expectedTeam.Members), len(team.Members))
				}
			}

			teamRepo.AssertExpectations(t)
		})
	}
}

func TestTeamService_GetTeam(t *testing.T) {
	tests := []struct {
		name          string
		teamName      string
		setupMocks    func(*MockTeamRepository)
		expectedError string
		expectedTeam  *entity.Team
	}{
		{
			name:     "successful get team",
			teamName: "team1",
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members: []entity.TeamMember{
						{UserID: "user1", Username: "user1", IsActive: true},
						{UserID: "user2", Username: "user2", IsActive: true},
					},
				}, nil)
			},
			expectedError: "",
			expectedTeam: &entity.Team{
				TeamName: "team1",
				Members: []entity.TeamMember{
					{UserID: "user1", Username: "user1", IsActive: true},
					{UserID: "user2", Username: "user2", IsActive: true},
				},
			},
		},
		{
			name:     "team not found",
			teamName: "team1",
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError: "NOT_FOUND",
			expectedTeam:  nil,
		},
		{
			name:     "get team error",
			teamName: "team1",
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(nil, errors.New("db error"))
			},
			expectedError: "db error",
			expectedTeam:  nil,
		},
		{
			name:     "team with no members",
			teamName: "team1",
			setupMocks: func(teamRepo *MockTeamRepository) {
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
			},
			expectedError: "",
			expectedTeam: &entity.Team{
				TeamName: "team1",
				Members:  []entity.TeamMember{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			teamRepo := new(MockTeamRepository)

			tt.setupMocks(teamRepo)

			svc := NewTeamService(teamRepo)
			ctx := context.Background()

			team, err := svc.GetTeam(ctx, tt.teamName)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, team)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, team)
				if tt.expectedTeam != nil {
					assert.Equal(t, tt.expectedTeam.TeamName, team.TeamName)
					assert.Equal(t, len(tt.expectedTeam.Members), len(team.Members))
				}
			}

			teamRepo.AssertExpectations(t)
		})
	}
}
