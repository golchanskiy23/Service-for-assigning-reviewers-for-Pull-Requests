package service

import (
	"context"

	"github.com/stretchr/testify/mock"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

type MockPullRequestRepository struct {
	mock.Mock
}

func (m *MockPullRequestRepository) CreatePR(ctx context.Context, pr *entity.PullRequest, reviewerIDs []string) error {
	args := m.Called(ctx, pr, reviewerIDs)
	return args.Error(0)
}

func (m *MockPullRequestRepository) GetPR(ctx context.Context, prID string) (*entity.PullRequest, error) {
	args := m.Called(ctx, prID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	pr, ok := args.Get(0).(*entity.PullRequest)
	if !ok {
		return nil, args.Error(1)
	}

	return pr, args.Error(1)
}

func (m *MockPullRequestRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	args := m.Called(ctx, prID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPullRequestRepository) UpdatePR(ctx context.Context, pr *entity.PullRequest) error {
	args := m.Called(ctx, pr)
	return args.Error(0)
}

func (m *MockPullRequestRepository) UpdateReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
	args := m.Called(ctx, prID, reviewerIDs)
	return args.Error(0)
}

func (m *MockPullRequestRepository) GetOpenPRsByReviewer(ctx context.Context, reviewerID string) ([]string, error) {
	args := m.Called(ctx, reviewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	pr, ok := args.Get(0).([]string)
	if !ok {
		return nil, args.Error(1)
	}

	return pr, args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) MassDeactivateAndReassign(ctx context.Context, teamName string, userIDs []string) error {
	args := m.Called(ctx, teamName, userIDs)
	return args.Error(0)
}

func (m *MockUserRepository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	pr, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, args.Error(1)
	}

	return pr, args.Error(1)
}

func (m *MockUserRepository) SetIsActive(ctx context.Context, userID string, active bool) error {
	args := m.Called(ctx, userID, active)
	return args.Error(0)
}

func (m *MockUserRepository) GetActiveUsersByTeam(
	ctx context.Context,
	teamName string,
	exclude []string,
) ([]*entity.User, error) {
	args := m.Called(ctx, teamName, exclude)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	users, ok := args.Get(0).([]*entity.User)
	if !ok {
		return nil, args.Error(1)
	}

	return users, args.Error(1)
}

func (m *MockUserRepository) GetPRsForReviewer(ctx context.Context, userID string) ([]*entity.PullRequestShort, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	prs, ok := args.Get(0).([]*entity.PullRequestShort)
	if !ok {
		return nil, args.Error(1)
	}

	return prs, args.Error(1)
}

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) AddTeam(ctx context.Context, team *entity.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamRepository) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	team, ok := args.Get(0).(*entity.Team)
	if !ok {
		return nil, args.Error(1)
	}

	return team, args.Error(1)
}

func (m *MockTeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	args := m.Called(ctx, teamName)
	return args.Bool(0), args.Error(1)
}
