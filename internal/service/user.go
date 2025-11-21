package service

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"context"
	"errors"
	"time"
)

type UserService struct {
	repo      postgres.UserRepository
	prRepo    postgres.PullRequestRepository
	teamRepo  postgres.TeamRepository
	prService *PRService
}

func NewUserService(repo postgres.UserRepository, prRepo postgres.PullRequestRepository, teamRepo postgres.TeamRepository, prService *PRService) *UserService {
	return &UserService{
		repo:      repo,
		prRepo:    prRepo,
		teamRepo:  teamRepo,
		prService: prService,
	}
}

func (s *UserService) ChangeActivateStatus(ctx context.Context, userID string, isActive bool) (*entity.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	user, err := s.repo.GetUser(queryCtx, userID)
	if err != nil {
		return nil, errors.New("NOT_FOUND")
	}

	if !isActive && user.IsActive {
		openPRs, err := s.prRepo.GetOpenPRsByReviewer(queryCtx, userID)
		if err != nil {
			return nil, err
		}

		maxPRs := 5
		if len(openPRs) > maxPRs {
			openPRs = openPRs[:maxPRs]
		}

		for _, prID := range openPRs {
			reassignCtx, reassignCancel := context.WithTimeout(queryCtx, 100*time.Millisecond)
			_, _, err := s.prService.ReassignReviewer(reassignCtx, prID, userID)
			reassignCancel()

			if err != nil {
				pr, err := s.prRepo.GetPR(queryCtx, prID)
				if err == nil {
					newReviewers := []string{}
					for _, reviewerID := range pr.AssignedReviewers {
						if reviewerID != userID {
							newReviewers = append(newReviewers, reviewerID)
						}
					}
					s.prRepo.UpdateReviewers(queryCtx, prID, newReviewers)
				}
			}
		}
	}

	err = s.repo.SetIsActive(queryCtx, userID, isActive)
	if err != nil {
		return nil, err
	}

	user, err = s.repo.GetUser(queryCtx, userID)
	if err != nil {
		return nil, errors.New("NOT_FOUND")
	}
	return user, nil
}

func (s *UserService) GetPRsAssignedTo(ctx context.Context, userID string) (string, []*entity.PullRequestShort, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
	defer cancel()

	_, err := s.repo.GetUser(queryCtx, userID)
	if err != nil {
		return "", nil, errors.New("NOT_FOUND")
	}

	prs, err := s.repo.GetPRsForReviewer(queryCtx, userID)
	if err != nil {
		return "", nil, err
	}

	return userID, prs, nil
}
