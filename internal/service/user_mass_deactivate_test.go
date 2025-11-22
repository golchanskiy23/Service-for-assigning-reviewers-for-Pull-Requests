package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

func TestUserService_MassDeactivate(t *testing.T) {
	ctx := t.Context()

	// set up mocks
	userRepo := new(MockUserRepository)
	prRepo := new(MockPullRequestRepository)
	teamRepo := new(MockTeamRepository)

	svc := NewUserService(userRepo, prRepo, teamRepo, nil)

	t.Run("flag true returns ONLY_DEACTIVATE", func(t *testing.T) {
		err := svc.MassDeactivate(ctx, []entity.User{{UserID: "u1"}}, true)
		assert.EqualError(t, err, "ONLY_DEACTIVATE")
	})

	t.Run("empty request returns EMPTY_REQUEST", func(t *testing.T) {
		err := svc.MassDeactivate(ctx, []entity.User{}, false)
		assert.EqualError(t, err, "EMPTY_REQUEST")
	})

	t.Run("different team in provided users returns DIFFERENT_TEAM", func(t *testing.T) {
		users := []entity.User{
			{UserID: "u1", TeamName: "team1"},
			{UserID: "u2", TeamName: "team2"},
		}
		err := svc.MassDeactivate(ctx, users, false)
		assert.EqualError(t, err, "DIFFERENT_TEAM")
	})

	t.Run("repo GetUser returns not found when team missing", func(t *testing.T) {
		users := []entity.User{{UserID: "u1", TeamName: ""}}
		userRepo.On("GetUser", mock.Anything, "u1").Return(nil, errors.New("NOT_FOUND"))
		defer userRepo.AssertExpectations(t)

		err := svc.MassDeactivate(ctx, users, false)
		assert.EqualError(t, err, "NOT_FOUND")
	})

	t.Run("successful mass deactivate calls repository with team and ids", func(t *testing.T) {
		users := []entity.User{{UserID: "u1", TeamName: "team1"}, {UserID: "u2", TeamName: "team1"}}
		userRepo.On("MassDeactivateAndReassign", mock.Anything, "team1", []string{"u1", "u2"}).Return(nil)
		defer userRepo.AssertExpectations(t)

		err := svc.MassDeactivate(ctx, users, false)
		assert.NoError(t, err)
	})
}
