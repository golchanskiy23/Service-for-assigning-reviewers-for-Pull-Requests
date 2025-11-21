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
	userQueryTimeout       = 300 * time.Millisecond
	userGetPRsQueryTimeout = 250 * time.Millisecond
	maxPRsToProcess        = 5
	reassignTimeout        = 100 * time.Millisecond
)

type UserService struct {
	repo      postgres.UserRepository
	prRepo    postgres.PullRequestRepository
	teamRepo  postgres.TeamRepository
	prService *PRService
}

//nolint:lll,revive // func
func NewUserService(repo postgres.UserRepository, prRepo postgres.PullRequestRepository, teamRepo postgres.TeamRepository, prService *PRService) *UserService {
	return &UserService{
		repo:      repo,
		prRepo:    prRepo,
		teamRepo:  teamRepo,
		prService: prService,
	}
}

//nolint:gocognit,nestif,revive,cyclop // Complex business logic for user activation status change
func (s *UserService) ChangeStatus(ctx context.Context, userID string, isActive bool) (*entity.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, userQueryTimeout)
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

		maxPRs := maxPRsToProcess
		if len(openPRs) > maxPRs {
			openPRs = openPRs[:maxPRs]
		}

		for _, prID := range openPRs {
			reassignCtx, reassignCancel := context.WithTimeout(queryCtx, reassignTimeout)
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

					if err := s.prRepo.UpdateReviewers(queryCtx, prID, newReviewers); err != nil {
						// Log error but continue processing other PRs.
						_ = err
					}
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

func (s *UserService) GetPRsAssignedTo(
	ctx context.Context,
	userID string,
) (string, []*entity.PullRequestShort, error) {
	queryCtx, cancel := context.WithTimeout(ctx, userGetPRsQueryTimeout)
	defer cancel()

	_, err := s.repo.GetUser(queryCtx, userID)
	if err != nil {
		const notFoundErr = "NOT_FOUND"
		return "", nil, errors.New(notFoundErr)
	}

	prs, err := s.repo.GetPRsForReviewer(queryCtx, userID)
	if err != nil {
		return "", nil, err
	}

	return userID, prs, nil
}
