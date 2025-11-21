package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

//nolint:maintidx // Complex test with many test cases
func TestUserService_ChangeActivateStatus(t *testing.T) {
	tests := []struct {
		setupMocks    func(*MockUserRepository, *MockPullRequestRepository, *MockTeamRepository, *PRService)
		expectedUser  *entity.User
		name          string
		userID        string
		expectedError string
		isActive      bool
	}{
		{
			name:     "successful activation",
			userID:   "user1",
			isActive: true,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil).Once()
				userRepo.On("SetIsActive", mock.Anything, "user1", true).Return(nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
			},
			expectedError: "",
			expectedUser: &entity.User{
				UserID:   "user1",
				Username: "testuser",
				TeamName: "team1",
				IsActive: true,
			},
		},
		{
			name:     "successful deactivation without open PRs",
			userID:   "user1",
			isActive: false,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				prRepo.On("GetOpenPRsByReviewer", mock.Anything, "user1").Return([]string{}, nil)
				userRepo.On("SetIsActive", mock.Anything, "user1", false).Return(nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil).Once()
			},
			expectedError: "",
			expectedUser: &entity.User{
				UserID:   "user1",
				Username: "testuser",
				TeamName: "team1",
				IsActive: false,
			},
		},
		{
			name:     "successful deactivation with open PRs - reassignment success",
			userID:   "user1",
			isActive: false,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				prRepo.On("GetOpenPRsByReviewer", mock.Anything, "user1").Return([]string{"pr1", "pr2"}, nil)
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR 1",
					AuthorID:          "user2",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user1"},
					CreatedAt:         &now,
				}, nil).Once()
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user2", "user1"}).Return([]*entity.User{
					{UserID: "user3", Username: "user3", TeamName: "team1", IsActive: true},
				}, nil)
				prRepo.On("UpdateReviewers", mock.Anything, "pr1", mock.Anything).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR 1",
					AuthorID:          "user2",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user3"},
					CreatedAt:         &now,
				}, nil).Once()
				prRepo.On("GetPR", mock.Anything, "pr2").Return(&entity.PullRequest{
					PullRequestID:     "pr2",
					PullRequestName:   "Test PR 2",
					AuthorID:          "user2",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user1"},
					CreatedAt:         &now,
				}, nil).Once()
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user2", "user1"}).Return([]*entity.User{
					{UserID: "user4", Username: "user4", TeamName: "team1", IsActive: true},
				}, nil)
				prRepo.On("UpdateReviewers", mock.Anything, "pr2", mock.Anything).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr2").Return(&entity.PullRequest{
					PullRequestID:     "pr2",
					PullRequestName:   "Test PR 2",
					AuthorID:          "user2",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user4"},
					CreatedAt:         &now,
				}, nil).Once()
				userRepo.On("SetIsActive", mock.Anything, "user1", false).Return(nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil).Once()
			},
			expectedError: "",
			expectedUser: &entity.User{
				UserID:   "user1",
				Username: "testuser",
				TeamName: "team1",
				IsActive: false,
			},
		},
		{
			name:     "successful deactivation with open PRs - reassignment failure, remove reviewer",
			userID:   "user1",
			isActive: false,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				prRepo.On("GetOpenPRsByReviewer", mock.Anything, "user1").Return([]string{"pr1"}, nil)
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR 1",
					AuthorID:          "user2",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user1"},
					CreatedAt:         &now,
				}, nil).Once()
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				userRepo.On(
					"GetActiveUsersByTeam",
					mock.Anything,
					"team1",
					[]string{"user2", "user1"},
				).Return([]*entity.User{}, nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR 1",
					AuthorID:          "user2",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user1"},
					CreatedAt:         &now,
				}, nil).Once()
				prRepo.On("UpdateReviewers", mock.Anything, "pr1", []string{}).Return(nil)
				userRepo.On("SetIsActive", mock.Anything, "user1", false).Return(nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil).Once()
			},
			expectedError: "",
			expectedUser: &entity.User{
				UserID:   "user1",
				Username: "testuser",
				TeamName: "team1",
				IsActive: false,
			},
		},
		{
			name:     "user not found",
			userID:   "user1",
			isActive: true,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError: "NOT_FOUND",
			expectedUser:  nil,
		},
		{
			name:     "get open PRs error",
			userID:   "user1",
			isActive: false,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				prRepo.On("GetOpenPRsByReviewer", mock.Anything, "user1").Return(nil, errors.New("db error"))
			},
			expectedError: "db error",
			expectedUser:  nil,
		},
		{
			name:     "set is active error",
			userID:   "user1",
			isActive: true,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil).Once()
				userRepo.On("SetIsActive", mock.Anything, "user1", true).Return(errors.New("db error"))
			},
			expectedError: "db error",
			expectedUser:  nil,
		},
		{
			name:     "get user after update error",
			userID:   "user1",
			isActive: true,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil).Once()
				userRepo.On("SetIsActive", mock.Anything, "user1", true).Return(nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(nil, errors.New("NOT_FOUND")).Once()
			},
			expectedError: "NOT_FOUND",
			expectedUser:  nil,
		},
		{
			name:     "deactivation with more than 5 open PRs",
			userID:   "user1",
			isActive: false,
			setupMocks: func(
				userRepo *MockUserRepository,
				prRepo *MockPullRequestRepository,
				teamRepo *MockTeamRepository,
				prService *PRService,
			) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil).Once()
				prRepo.On(
					"GetOpenPRsByReviewer",
					mock.Anything,
					"user1",
				).Return(
					[]string{"pr1", "pr2", "pr3", "pr4", "pr5", "pr6", "pr7"},
					nil,
				)
				now := time.Now()
				prIDs := []string{"pr1", "pr2", "pr3", "pr4", "pr5"}
				for _, prID := range prIDs {
					prRepo.On("GetPR", mock.Anything, prID).Return(&entity.PullRequest{
						PullRequestID:     prID,
						PullRequestName:   "Test PR",
						AuthorID:          "user2",
						Status:            entity.OPEN,
						AssignedReviewers: []string{"user1"},
						CreatedAt:         &now,
					}, nil).Once()
					userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
						UserID:   "user1",
						Username: "testuser",
						TeamName: "team1",
						IsActive: true,
					}, nil).Once()
					userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user2", "user1"}).Return([]*entity.User{
						{UserID: "user3", Username: "user3", TeamName: "team1", IsActive: true},
					}, nil).Once()
					prRepo.On("UpdateReviewers", mock.Anything, prID, mock.Anything).Return(nil).Once()
					prRepo.On("GetPR", mock.Anything, prID).Return(&entity.PullRequest{
						PullRequestID:     prID,
						PullRequestName:   "Test PR",
						AuthorID:          "user2",
						Status:            entity.OPEN,
						AssignedReviewers: []string{"user3"},
						CreatedAt:         &now,
					}, nil).Once()
				}
				userRepo.On("SetIsActive", mock.Anything, "user1", false).Return(nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: false,
				}, nil).Once()
			},
			expectedError: "",
			expectedUser: &entity.User{
				UserID:   "user1",
				Username: "testuser",
				TeamName: "team1",
				IsActive: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			prRepo := new(MockPullRequestRepository)
			teamRepo := new(MockTeamRepository)
			prService := NewPRService(prRepo, userRepo, teamRepo)

			tt.setupMocks(userRepo, prRepo, teamRepo, prService)

			svc := NewUserService(userRepo, prRepo, teamRepo, prService)
			ctx := t.Context()

			user, err := svc.ChangeStatus(ctx, tt.userID, tt.isActive)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)

				if tt.expectedUser != nil {
					assert.Equal(t, tt.expectedUser.UserID, user.UserID)
					assert.Equal(t, tt.expectedUser.IsActive, user.IsActive)
				}
			}

			userRepo.AssertExpectations(t)
			prRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetPRsAssignedTo(t *testing.T) {
	tests := []struct {
		setupMocks     func(*MockUserRepository)
		name           string
		userID         string
		expectedError  string
		expectedUserID string
		expectedPRs    int
	}{
		{
			name:   "successful get PRs",
			userID: "user1",
			setupMocks: func(userRepo *MockUserRepository) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				userRepo.On("GetPRsForReviewer", mock.Anything, "user1").Return([]*entity.PullRequestShort{
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
			expectedError:  "",
			expectedPRs:    2,
			expectedUserID: "user1",
		},
		{
			name:   "user not found",
			userID: "user1",
			setupMocks: func(userRepo *MockUserRepository) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError:  "NOT_FOUND",
			expectedPRs:    0,
			expectedUserID: "",
		},
		{
			name:   "get PRs error",
			userID: "user1",
			setupMocks: func(userRepo *MockUserRepository) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				userRepo.On("GetPRsForReviewer", mock.Anything, "user1").Return(nil, errors.New("db error"))
			},
			expectedError:  "db error",
			expectedPRs:    0,
			expectedUserID: "",
		},
		{
			name:   "no PRs assigned",
			userID: "user1",
			setupMocks: func(userRepo *MockUserRepository) {
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				userRepo.On("GetPRsForReviewer", mock.Anything, "user1").Return([]*entity.PullRequestShort{}, nil)
			},
			expectedError:  "",
			expectedPRs:    0,
			expectedUserID: "user1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			prRepo := new(MockPullRequestRepository)
			teamRepo := new(MockTeamRepository)
			prService := NewPRService(prRepo, userRepo, teamRepo)

			tt.setupMocks(userRepo)

			svc := NewUserService(userRepo, prRepo, teamRepo, prService)
			ctx := t.Context()

			userID, prs, err := svc.GetPRsAssignedTo(ctx, tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, prs)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, userID)
				assert.NotNil(t, prs)
				assert.Len(t, prs, tt.expectedPRs)
			}

			userRepo.AssertExpectations(t)
		})
	}
}
