package handlers

import (
	"log/slog"

	//nolint:revive // dependency
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/service"
)

type Services struct {
	Log          *slog.Logger
	TeamService  TeamServiceInterface
	UserService  UserServiceInterface
	PRService    PRServiceInterface
	LoadService  LoadServiceInterface
	StatsService StatsServiceInterface
}

//nolint:revive // long line
func CreateNewService(repo *postgres.Repository, logger *slog.Logger) *Services {
	prService := service.NewPRService(repo.PullRequests, repo.Users, repo.Teams)

	return &Services{
		Log:         logger,
		TeamService: service.NewTeamService(repo.Teams),
		UserService: service.NewUserService(
			repo.Users,
			repo.PullRequests,
			repo.Teams,
			prService,
		),
		PRService:    prService,
		LoadService:  &service.LoadService{},
		StatsService: service.NewStatsService(repo.Stats),
	}
}
