package handlers

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/service"
)

type Services struct {
	TeamService *service.TeamService
	UserService *service.UserService
	PRService   *service.PRService
}

func CreateNewService(repo *postgres.Repository) *Services {
	return &Services{
		TeamService: service.NewTeamService(repo.Teams),
		UserService: service.NewUserService(repo.Users),
		PRService:   service.NewPRService(repo.PullRequests, repo.Users, repo.Teams),
	}
}
