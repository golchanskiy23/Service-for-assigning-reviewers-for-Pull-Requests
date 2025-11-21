package service

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"context"
	"errors"
	"time"
)

type TeamService struct {
	repo postgres.TeamRepository
}

func NewTeamService(repo postgres.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) AddTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	exists, err := s.repo.TeamExists(queryCtx, team.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("TEAM_EXISTS")
	}

	err = s.repo.AddTeam(queryCtx, team)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			return nil, err
		}
		return nil, err
	}

	return s.repo.GetTeam(queryCtx, team.TeamName)
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
	defer cancel()

	team, err := s.repo.GetTeam(queryCtx, teamName)
	if err != nil {
		return nil, err
	}
	return team, nil
}
