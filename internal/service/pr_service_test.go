package service

/*import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)*/

/*
func TestPRService_CreatePR(t *testing.T) {
	tests := []struct {
		name          string
		prID          string
		prName        string
		authorID      string
		setupMocks    func(*MockPullRequestRepository, *MockUserRepository, *MockTeamRepository)
		expectedError string
		expectedPR    *entity.PullRequest
		expectedMsg   string
	}{
		{
			name:     "successful PR creation with reviewers",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(false, nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user1"}).Return([]*entity.User{
					{UserID: "user2", Username: "user2", TeamName: "team1", IsActive: true},
					{UserID: "user3", Username: "user3", TeamName: "team1", IsActive: true},
				}, nil)
				now := time.Now()
				prRepo.On("CreatePR", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2", "user3"},
					CreatedAt:         &now,
				}, nil)
			},
			expectedError: "",
			expectedMsg:   "",
		},
		{
			name:     "PR already exists",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(true, nil)
			},
			expectedError: "PR_EXISTS",
			expectedMsg:   "",
		},
		{
			name:     "PR exists check error",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(false, errors.New("db error"))
			},
			expectedError: "db error",
			expectedMsg:   "",
		},
		{
			name:     "author not found",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(false, nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError: "NOT_FOUND",
			expectedMsg:   "",
		},
		{
			name:     "team not found",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(false, nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError: "NOT_FOUND",
			expectedMsg:   "",
		},
		{
			name:     "no candidates for reviewers",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(false, nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user1"}).Return([]*entity.User{}, nil)
				now := time.Now()
				prRepo.On("CreatePR", mock.Anything, mock.Anything, []string{}).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{},
					CreatedAt:         &now,
				}, nil)
			},
			expectedError: "",
			expectedMsg:   "",
		},
		{
			name:     "create PR error",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(false, nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user1"}).Return([]*entity.User{
					{UserID: "user2", Username: "user2", TeamName: "team1", IsActive: true},
				}, nil)
				prRepo.On("CreatePR", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("PR_EXISTS"))
			},
			expectedError: "PR_EXISTS",
			expectedMsg:   "",
		},
		{
			name:     "get PR after creation error",
			prID:     "pr1",
			prName:   "Test PR",
			authorID: "user1",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository, teamRepo *MockTeamRepository) {
				prRepo.On("PRExists", mock.Anything, "pr1").Return(false, nil)
				userRepo.On("GetUser", mock.Anything, "user1").Return(&entity.User{
					UserID:   "user1",
					Username: "testuser",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				teamRepo.On("GetTeam", mock.Anything, "team1").Return(&entity.Team{
					TeamName: "team1",
					Members:  []entity.TeamMember{},
				}, nil)
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user1"}).Return([]*entity.User{
					{UserID: "user2", Username: "user2", TeamName: "team1", IsActive: true},
				}, nil)
				prRepo.On("CreatePR", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(nil, errors.New("db error"))
			},
			expectedError: "db error",
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prRepo := new(MockPullRequestRepository)
			userRepo := new(MockUserRepository)
			teamRepo := new(MockTeamRepository)

			tt.setupMocks(prRepo, userRepo, teamRepo)

			service := NewPRService(prRepo, userRepo, teamRepo)
			ctx := t.Context()

			pr, msg, err := service.CreatePR(ctx, tt.prID, tt.prName, tt.authorID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, pr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pr)
				assert.Equal(t, tt.prID, pr.PullRequestID)
				assert.Equal(t, tt.prName, pr.PullRequestName)
				assert.Equal(t, tt.authorID, pr.AuthorID)
			}

			assert.Equal(t, tt.expectedMsg, msg)

			prRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			teamRepo.AssertExpectations(t)
		})
	}
}

func TestPRService_MergePR(t *testing.T) {
	tests := []struct {
		name           string
		prID           string
		setupMocks     func(*MockPullRequestRepository)
		expectedError  string
		expectedStatus entity.PRStatus
	}{
		{
			name: "successful merge",
			prID: "pr1",
			setupMocks: func(prRepo *MockPullRequestRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
				}, nil).Once()
				prRepo.On("UpdatePR", mock.Anything, mock.MatchedBy(func(pr *entity.PullRequest) bool {
					return pr.Status == entity.MERGED
				})).Return(nil)
				mergedTime := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.MERGED,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
					MergedAt:          &mergedTime,
				}, nil).Once()
			},
			expectedError:  "",
			expectedStatus: entity.MERGED,
		},
		{
			name: "PR not found",
			prID: "pr1",
			setupMocks: func(prRepo *MockPullRequestRepository) {
				prRepo.On("GetPR", mock.Anything, "pr1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError:  "NOT_FOUND",
			expectedStatus: "",
		},
		{
			name: "PR already merged",
			prID: "pr1",
			setupMocks: func(prRepo *MockPullRequestRepository) {
				now := time.Now()
				mergedTime := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.MERGED,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
					MergedAt:          &mergedTime,
				}, nil)
			},
			expectedError:  "",
			expectedStatus: entity.MERGED,
		},
		{
			name: "update PR error",
			prID: "pr1",
			setupMocks: func(prRepo *MockPullRequestRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
				}, nil).Once()
				prRepo.On("UpdatePR", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedError:  "db error",
			expectedStatus: "",
		},
		{
			name: "get PR after update error",
			prID: "pr1",
			setupMocks: func(prRepo *MockPullRequestRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
				}, nil).Once()
				prRepo.On("UpdatePR", mock.Anything, mock.Anything).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(nil, errors.New("db error")).Once()
			},
			expectedError:  "db error",
			expectedStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prRepo := new(MockPullRequestRepository)
			userRepo := new(MockUserRepository)
			teamRepo := new(MockTeamRepository)

			tt.setupMocks(prRepo)

			service := NewPRService(prRepo, userRepo, teamRepo)
			ctx := t.Context()

			pr, err := service.MergePR(ctx, tt.prID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, pr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pr)

				if tt.expectedStatus != "" {
					assert.Equal(t, tt.expectedStatus, pr.Status)
				}
			}

			prRepo.AssertExpectations(t)
		})
	}
}

func TestPRService_ReassignReviewer(t *testing.T) {
	tests := []struct {
		name          string
		prID          string
		oldReviewerID string
		setupMocks    func(*MockPullRequestRepository, *MockUserRepository)
		expectedError string
		expectedPR    bool
		expectedNewID bool
	}{
		{
			name:          "successful reassignment",
			prID:          "pr1",
			oldReviewerID: "user2",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2", "user3"},
					CreatedAt:         &now,
				}, nil).Once()
				userRepo.On("GetUser", mock.Anything, "user2").Return(&entity.User{
					UserID:   "user2",
					Username: "user2",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user1", "user2"}).Return([]*entity.User{
					{UserID: "user4", Username: "user4", TeamName: "team1", IsActive: true},
				}, nil)
				prRepo.On("UpdateReviewers", mock.Anything, "pr1", mock.Anything).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user4", "user3"},
					CreatedAt:         &now,
				}, nil).Once()
			},
			expectedError: "",
			expectedPR:    true,
			expectedNewID: true,
		},
		{
			name:          "PR not found",
			prID:          "pr1",
			oldReviewerID: "user2",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				prRepo.On("GetPR", mock.Anything, "pr1").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError: "NOT_FOUND",
			expectedPR:    false,
			expectedNewID: false,
		},
		{
			name:          "PR already merged",
			prID:          "pr1",
			oldReviewerID: "user2",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				now := time.Now()
				mergedTime := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.MERGED,
					AssignedReviewers: []string{"user2"},
					CreatedAt:         &now,
					MergedAt:          &mergedTime,
				}, nil)
			},
			expectedError: "PR_MERGED",
			expectedPR:    false,
			expectedNewID: false,
		},
		{
			name:          "reviewer not assigned",
			prID:          "pr1",
			oldReviewerID: "user5",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2", "user3"},
					CreatedAt:         &now,
				}, nil)
			},
			expectedError: "NOT_ASSIGNED",
			expectedPR:    false,
			expectedNewID: false,
		},
		{
			name:          "old reviewer not found",
			prID:          "pr1",
			oldReviewerID: "user2",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2", "user3"},
					CreatedAt:         &now,
				}, nil)
				userRepo.On("GetUser", mock.Anything, "user2").Return(nil, errors.New("NOT_FOUND"))
			},
			expectedError: "NOT_FOUND",
			expectedPR:    false,
			expectedNewID: false,
		},
		{
			name:          "no candidate available",
			prID:          "pr1",
			oldReviewerID: "user2",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2", "user3"},
					CreatedAt:         &now,
				}, nil)
				userRepo.On("GetUser", mock.Anything, "user2").Return(&entity.User{
					UserID:   "user2",
					Username: "user2",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				userRepo.On(
					"GetActiveUsersByTeam",
					mock.Anything,
					"team1",
					[]string{"user1", "user2"},
				).Return([]*entity.User{}, nil)
			},
			expectedError: "NO_CANDIDATE",
			expectedPR:    false,
			expectedNewID: false,
		},
		{
			name:          "update reviewers error",
			prID:          "pr1",
			oldReviewerID: "user2",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2", "user3"},
					CreatedAt:         &now,
				}, nil).Once()
				userRepo.On("GetUser", mock.Anything, "user2").Return(&entity.User{
					UserID:   "user2",
					Username: "user2",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user1", "user2"}).Return([]*entity.User{
					{UserID: "user4", Username: "user4", TeamName: "team1", IsActive: true},
				}, nil)
				prRepo.On("UpdateReviewers", mock.Anything, "pr1", mock.Anything).Return(errors.New("db error"))
			},
			expectedError: "db error",
			expectedPR:    false,
			expectedNewID: false,
		},
		{
			name:          "get PR after update error",
			prID:          "pr1",
			oldReviewerID: "user2",
			setupMocks: func(prRepo *MockPullRequestRepository, userRepo *MockUserRepository) {
				now := time.Now()
				prRepo.On("GetPR", mock.Anything, "pr1").Return(&entity.PullRequest{
					PullRequestID:     "pr1",
					PullRequestName:   "Test PR",
					AuthorID:          "user1",
					Status:            entity.OPEN,
					AssignedReviewers: []string{"user2", "user3"},
					CreatedAt:         &now,
				}, nil).Once()
				userRepo.On("GetUser", mock.Anything, "user2").Return(&entity.User{
					UserID:   "user2",
					Username: "user2",
					TeamName: "team1",
					IsActive: true,
				}, nil)
				userRepo.On("GetActiveUsersByTeam", mock.Anything, "team1", []string{"user1", "user2"}).Return([]*entity.User{
					{UserID: "user4", Username: "user4", TeamName: "team1", IsActive: true},
				}, nil)
				prRepo.On("UpdateReviewers", mock.Anything, "pr1", mock.Anything).Return(nil)
				prRepo.On("GetPR", mock.Anything, "pr1").Return(nil, errors.New("db error")).Once()
			},
			expectedError: "db error",
			expectedPR:    false,
			expectedNewID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prRepo := new(MockPullRequestRepository)
			userRepo := new(MockUserRepository)
			teamRepo := new(MockTeamRepository)

			tt.setupMocks(prRepo, userRepo)

			service := NewPRService(prRepo, userRepo, teamRepo)
			ctx := t.Context()

			pr, newID, err := service.ReassignReviewer(ctx, tt.prID, tt.oldReviewerID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, pr)
				assert.Empty(t, newID)
			} else {
				assert.NoError(t, err)

				if tt.expectedPR {
					assert.NotNil(t, pr)
				}

				if tt.expectedNewID {
					assert.NotEmpty(t, newID)
				}
			}

			prRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}*/
