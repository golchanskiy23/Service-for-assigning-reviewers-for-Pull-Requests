package service

import "errors"

type TeamService struct {
	repo TeamRepository
}

func NewTeamService(r TeamRepository) *TeamService {
	return &TeamService{repo: r}
}

func (s *TeamService) AddTeam(name string, users []entity.User) (*entity.Team, error) {
	// доменная валидация
	if name == "" {
		return nil, errors.New("team name required")
	}

	return s.repo.CreateTeam(name, users)
}

func (s *TeamService) GetTeam(name string) (*entity.Team, error) {
	return s.repo.GetTeam(name)
}
