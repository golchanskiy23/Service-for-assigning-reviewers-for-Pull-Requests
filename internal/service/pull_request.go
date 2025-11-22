package service

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	//nolint:revive // necessary import
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
)

const (
	prQueryTimeout = 250 * time.Millisecond
	emptyString    = ""
	notFoundErr    = "NOT_FOUND"
	zeroLength     = 0
)

type PRService struct {
	repo     postgres.PullRequestRepository
	userRepo postgres.UserRepository
	teamRepo postgres.TeamRepository
}

//nolint:revive // func
func NewPRService(r postgres.PullRequestRepository, u postgres.UserRepository, t postgres.TeamRepository) *PRService {
	return &PRService{
		repo:     r,
		userRepo: u,
		teamRepo: t,
	}
}

//nolint:revive,cyclop // Complex business logic for PR creation
func (s *PRService) CreatePR(
	ctx context.Context,
	prID, prName, authorID string,
) (*entity.PullRequest, string, error) {
	queryCtx, cancel := context.WithTimeout(ctx, prQueryTimeout)

	defer cancel()

	exists, err := s.repo.PRExists(queryCtx, prID)
	if err != nil {
		return nil, emptyString, err
	}

	if exists {
		return nil, emptyString, entity.ErrPRExists
	}

	author, err := s.userRepo.GetUser(queryCtx, authorID)
	if err != nil {
		return nil, emptyString, entity.ErrNotFound
	}

	_, err = s.teamRepo.GetTeam(queryCtx, author.TeamName)
	if err != nil {
		return nil, emptyString, entity.ErrNotFound
	}

	candidates, err := s.userRepo.GetActiveUsersByTeam(queryCtx, author.TeamName, []string{authorID})
	if err != nil {
		return nil, emptyString, err
	}

	reviewerIDs := []string{}

	if len(candidates) > 0 {
		shuffled := make([]*entity.User, len(candidates))
		copy(shuffled, candidates)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		count := 2
		if len(shuffled) < count {
			count = len(shuffled)
		}

		for i := range count {
			reviewerIDs = append(reviewerIDs, shuffled[i].UserID)
		}
	}

	now := time.Now()
	pr := &entity.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            entity.OPEN,
		AssignedReviewers: reviewerIDs,
		CreatedAt:         &now,
	}

	err = s.repo.CreatePR(queryCtx, pr, reviewerIDs)
	if err != nil {
		if errors.Is(err, entity.ErrPRExists) {
			return nil, emptyString, err
		}

		return nil, emptyString, err
	}

	createdPR, err := s.repo.GetPR(queryCtx, prID)
	if err != nil {
		return nil, emptyString, err
	}

	return createdPR, emptyString, nil
}

//nolint:revive // func
func (s *PRService) MergePR(ctx context.Context, prID string) (*entity.PullRequest, error) {
	queryCtx, cancel := context.WithTimeout(ctx, prQueryTimeout)
	defer cancel()

	pr, err := s.repo.GetPR(queryCtx, prID)
	if err != nil {
		return nil, entity.ErrNotFound
	}

	if pr.Status == entity.MERGED {
		return pr, nil
	}

	pr.Status = entity.MERGED
	now := time.Now()
	pr.MergedAt = &now

	err = s.repo.UpdatePR(queryCtx, pr)
	if err != nil {
		return nil, err
	}

	return s.repo.GetPR(queryCtx, prID)
}

//nolint:revive,cyclop // Complex business logic for PR reassignment
func (s *PRService) ReassignReviewer(
	ctx context.Context,
	prID, oldReviewerID string,
) (*entity.PullRequest, string, error) {
	queryCtx, cancel := context.WithTimeout(ctx, prQueryTimeout)

	defer cancel()

	pr, e := s.repo.GetPR(queryCtx, prID)
	if e != nil {
		return nil, emptyString, entity.ErrNotFound
	}

	if pr.Status == entity.MERGED {
		return nil, emptyString, entity.ErrPRMerged
	}

	// If oldReviewerID is empty, interpret as "assign a new reviewer" (append)
	if oldReviewerID == emptyString {
		// use PR author team to find candidates (same as create PR flow)
		author, err := s.userRepo.GetUser(queryCtx, pr.AuthorID)
		if err != nil {
			return nil, emptyString, entity.ErrNotFound
		}

		// exclude author and any already assigned reviewers
		exclude := make([]string, 0, len(pr.AssignedReviewers)+1)
		exclude = append(exclude, pr.AuthorID)
		exclude = append(exclude, pr.AssignedReviewers...)

		candidates, err := s.userRepo.GetActiveUsersByTeam(queryCtx, author.TeamName, exclude)
		if err != nil {
			return nil, emptyString, err
		}
		// If PR currently has no reviewers, allow assigning up to 2 candidates (0..2)
		if len(pr.AssignedReviewers) == 0 {
			// if no candidates found, it's acceptable — return current PR without error
			if len(candidates) == zeroLength {
				return pr, emptyString, nil
			}

			// shuffle candidates and pick up to 2
			shuffled := make([]*entity.User, len(candidates))
			copy(shuffled, candidates)
			rand.Shuffle(len(shuffled), func(i, j int) {
				shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
			})

			pick := 2
			if len(shuffled) < pick {
				pick = len(shuffled)
			}

			selected := make([]string, 0, pick)
			for i := range pick {
				selected = append(selected, shuffled[i].UserID)
			}

			err = s.repo.UpdateReviewers(queryCtx, prID, selected)
			if err != nil {
				return nil, emptyString, err
			}

			updatedPR, error1 := s.repo.GetPR(queryCtx, prID)
			if error1 != nil {
				return nil, emptyString, err
			}

			// return comma-separated list of assigned IDs (may be 1 or 2)
			return updatedPR, strings.Join(selected, ","), nil
		}

		// otherwise (existing reviewers present) pick a single reviewer to append
		if len(candidates) == zeroLength {
			return nil, emptyString, entity.ErrNoCandidate
		}

		//nolint:gosec // math/rand is sufficient for selecting a random reviewer
		newReviewer := candidates[rand.Intn(len(candidates))]

		newReviewers := make([]string, 0, len(pr.AssignedReviewers)+1)
		newReviewers = append(newReviewers, pr.AssignedReviewers...)
		newReviewers = append(newReviewers, newReviewer.UserID)

		err = s.repo.UpdateReviewers(queryCtx, prID, newReviewers)
		if err != nil {
			return nil, emptyString, err
		}

		updatedPR, err := s.repo.GetPR(queryCtx, prID)
		if err != nil {
			return nil, emptyString, err
		}

		return updatedPR, newReviewer.UserID, nil
	}

	// Otherwise perform replacement of the specified old reviewer
	found := false

	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldReviewerID {
			found = true
			break
		}
	}

	if !found {
		return nil, emptyString, entity.ErrNotAssigned
	}

	oldReviewer, err := s.userRepo.GetUser(queryCtx, oldReviewerID)
	if err != nil {
		return nil, emptyString, entity.ErrNotFound
	}

	exclude := []string{pr.AuthorID, oldReviewerID}

	candidates, err := s.userRepo.GetActiveUsersByTeam(queryCtx, oldReviewer.TeamName, exclude)
	if err != nil {
		return nil, emptyString, err
	}

	if len(candidates) == zeroLength {
		return nil, emptyString, entity.ErrNoCandidate
	}

	//nolint:gosec // math/rand is sufficient for selecting a random reviewer
	newReviewer := candidates[rand.Intn(len(candidates))]

	newReviewers := make([]string, len(pr.AssignedReviewers))
	copy(newReviewers, pr.AssignedReviewers)

	for i, reviewerID := range newReviewers {
		if reviewerID == oldReviewerID {
			newReviewers[i] = newReviewer.UserID
			break
		}
	}

	err = s.repo.UpdateReviewers(queryCtx, prID, newReviewers)
	if err != nil {
		return nil, emptyString, err
	}

	updatedPR, err := s.repo.GetPR(queryCtx, prID)
	if err != nil {
		return nil, emptyString, err
	}

	return updatedPR, newReviewer.UserID, nil
}
