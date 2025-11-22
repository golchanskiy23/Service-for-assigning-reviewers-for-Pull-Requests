package service

import (
	"context"

	//nolint:revive // import
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
)

type StatsService struct {
	statsRepo postgres.StatsRepository
}

func NewStatsService(s postgres.StatsRepository) *StatsService {
	return &StatsService{statsRepo: s}
}

//nolint:revive // monolith func
func (s *StatsService) GetAssignedCountPerPR(ctx context.Context) (map[string]int, error) {
	return s.statsRepo.GetAssignedReviewersCountPerPR(ctx)
}

//nolint:revive // monolith func
func (s *StatsService) GetOpenPRCountPerUser(ctx context.Context) (map[string]int, error) {
	return s.statsRepo.GetOpenPRCountPerUser(ctx)
}
