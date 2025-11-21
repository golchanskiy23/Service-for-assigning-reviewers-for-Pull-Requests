package service

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"context"
)

type TeamService struct {
	repo postgres.TeamRepository
}

func NewTeamService(repo postgres.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string) error {
	return s.repo.Add(ctx, nil)
}

func (s *TeamService) Get(ctx context.Context, id int) (*postgres.Team, error) {
	return s.repo.Get(ctx, id)
}
