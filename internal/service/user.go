package service

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"context"
)

type UserService struct {
	repo postgres.UserRepository
}

func NewUserService(repo postgres.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ActivateUser(ctx context.Context, id int) error {
	return s.repo.SetIsActive(ctx, id, true)
}

func (s *UserService) GetUserReviews(ctx context.Context, id int) ([]string, error) {
	return s.repo.GetReview(ctx, id)
}

/*
func (s *UserService) SetUserActive(userID int64, active bool) (*entity.User, error) {
	return s.repo.UpdateActive(userID, active)
}

func (s *UserService) GetPRsAssignedTo(userID int64) ([]entity.PullRequest, error) {
	return s.repo.GetPRsForReviewer(userID)
}
*/
