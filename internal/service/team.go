package service

import (
	"context"
	"errors"
	"time"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	//nolint:revive // necessary import
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
)

const (
	teamQueryTimeout    = 300 * time.Millisecond
	teamGetQueryTimeout = 250 * time.Millisecond
)

type TeamService struct {
	repo postgres.TeamRepository
}

func NewTeamService(repo postgres.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

//nolint:revive // func
func (s *TeamService) AddTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	queryCtx, cancel := context.WithTimeout(ctx, teamQueryTimeout)
	defer cancel()

	exists, err := s.repo.TeamExists(queryCtx, team.TeamName)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, entity.ErrTeamExists
	}

	err = s.repo.AddTeam(queryCtx, team)
	if err != nil {
		if errors.Is(err, entity.ErrTeamExists) {
			return nil, err
		}

		return nil, err
	}

	return s.repo.GetTeam(queryCtx, team.TeamName)
}

//nolint:revive // func
func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	queryCtx, cancel := context.WithTimeout(ctx, teamGetQueryTimeout)
	defer cancel()

	team, err := s.repo.GetTeam(queryCtx, teamName)
	if err != nil {
		return nil, err
	}

	return team, nil
}
