package handlers

import (
	"context"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
)

type PRServiceInterface interface {
	CreatePR(
		ctx context.Context,
		prID, prName, authorID string,
	) (
		*entity.PullRequest,
		string,
		error,
	)
	MergePR(ctx context.Context, prID string) (*entity.PullRequest, error)
	ReassignReviewer(
		ctx context.Context,
		prID, oldReviewerID string,
	) (
		*entity.PullRequest,
		string,
		error,
	)
}

type UserServiceInterface interface {
	ChangeStatus(
		ctx context.Context,
		userID string,
		isActive bool,
	) (*entity.User, error)
	GetPRsAssignedTo(
		ctx context.Context,
		userID string,
	) (
		string,
		[]*entity.PullRequestShort,
		error,
	)
	MassDeactivate(ctx context.Context, users []entity.User, flag bool) error
}

type TeamServiceInterface interface {
	AddTeam(ctx context.Context, team *entity.Team) (*entity.Team, error)
	GetTeam(ctx context.Context, teamName string) (*entity.Team, error)
}

type LoadServiceInterface interface {
	RunLoadTest(rate vegeta.Rate, duration time.Duration)
}

type StatsServiceInterface interface {
	GetAssignedCountPerPR(ctx context.Context) (map[string]int, error)
	GetOpenPRCountPerUser(ctx context.Context) (map[string]int, error)
}
