package handlers

import "Service-for-assigning-reviewers-for-Pull-Requests/internal/service"

type ServiceExecution struct {
	TeamService *service.TeamService
	UserService *service.UserService
	PrService   *service.PRService
}
