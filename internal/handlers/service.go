package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/service"
	"context"
)

type PRServiceInterface interface {
	CreatePR(ctx context.Context, prID, prName, authorID string) (*entity.PullRequest, string, error)
	MergePR(ctx context.Context, prID string) (*entity.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*entity.PullRequest, string, error)
}

type UserServiceInterface interface {
	ChangeActivateStatus(ctx context.Context, userID string, isActive bool) (*entity.User, error)
	GetPRsAssignedTo(ctx context.Context, userID string) (string, []*entity.PullRequestShort, error)
}

type TeamServiceInterface interface {
	AddTeam(ctx context.Context, team *entity.Team) (*entity.Team, error)
	GetTeam(ctx context.Context, teamName string) (*entity.Team, error)
}

type Services struct {
	TeamService TeamServiceInterface
	UserService UserServiceInterface
	PRService   PRServiceInterface
}

func CreateNewService(repo *postgres.Repository) *Services {
	prService := service.NewPRService(repo.PullRequests, repo.Users, repo.Teams)
	return &Services{
		TeamService: service.NewTeamService(repo.Teams),
		UserService: service.NewUserService(repo.Users, repo.PullRequests, repo.Teams, prService),
		PRService:   prService,
	}
}
