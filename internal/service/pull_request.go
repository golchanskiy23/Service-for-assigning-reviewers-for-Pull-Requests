package service

import "Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"

type PRService struct {
	repo     postgres.PullRequestRepository
	userRepo postgres.UserRepository
	teamRepo postgres.TeamRepository
}

func NewPRService(r postgres.PullRequestRepository, u postgres.UserRepository, t postgres.TeamRepository) *PRService {
	return &PRService{
		repo:     r,
		userRepo: u,
		teamRepo: t,
	}
}

/*
func (s *PRService) CreatePR(pr *entity.PullRequest) (*entity.PullRequest, error) {
	team, err := s.teamRepo.GetTeam(pr.TeamName)
	if err != nil {
		return nil, errors.New("team not found")
	}

	// Получаем активных ревьюверов
	exclude := []int64{pr.AuthorID}
	reviewers, err := s.repo.GetActiveReviewers(team.Name, exclude)
	if err != nil {
		return nil, err
	}

	// Перемешиваем
	rand.Shuffle(len(reviewers), func(i, j int) {
		reviewers[i], reviewers[j] = reviewers[j], reviewers[i]
	})

	// Берем до 2 ревьюверов
	for i := 0; i < len(reviewers) && len(pr.Reviewers) < 2; i++ {
		pr.Reviewers = append(pr.Reviewers, reviewers[i].ID)
	}

	pr.Status = entity.PROpen

	return pr, s.repo.CreatePR(pr)
}

func (s *PRService) MergePR(id int64) (*entity.PullRequest, error) {
	pr, err := s.repo.GetPR(id)
	if err != nil {
		return nil, errors.New("not found")
	}

	if pr.Status == entity.PRMerged {
		return pr, nil // идемпотентность
	}

	pr.Status = entity.PRMerged
	return pr, s.repo.UpdatePR(pr)
}

func (s *PRService) ReassignReviewer(prID, oldReviewerID int64) (*entity.PullRequest, int64, error) {
	pr, err := s.repo.GetPR(prID)
	if err != nil {
		return nil, 0, errors.New("pr not found")
	}

	if pr.Status == entity.PRMerged {
		return nil, 0, errors.New("PR_MERGED")
	}

	// Проверяем, что юзер был ревьювером
	found := false
	for _, r := range pr.Reviewers {
		if r == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, 0, errors.New("NOT_ASSIGNED")
	}

	// Нужно получить команду ревьювера
	reviewer, err := s.userRepo.GetUser(oldReviewerID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	team, err := s.teamRepo.GetTeam(reviewer.TeamName)
	if err != nil {
		return nil, 0, errors.New("team not found")
	}

	exclude := append([]int64{pr.AuthorID, oldReviewerID}, pr.Reviewers...)
	candidates, err := s.repo.GetActiveReviewers(team.Name, exclude)
	if err != nil {
		return nil, 0, err
	}

	if len(candidates) == 0 {
		return nil, 0, errors.New("NO_CANDIDATE")
	}

	newReviewer := candidates[rand.Intn(len(candidates))]

	// заменить
	for i, r := range pr.Reviewers {
		if r == oldReviewerID {
			pr.Reviewers[i] = newReviewer.ID
			break
		}
	}

	return pr, newReviewer.ID, s.repo.UpdatePR(pr)
}
*/
