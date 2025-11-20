package service

type TeamRepository interface {
	CreateTeam(name string, users []entity.User) (*entity.Team, error)
	GetTeam(name string) (*entity.Team, error)
}

type UserRepository interface {
	UpdateActive(userID int64, active bool) (*entity.User, error)
	GetUser(userID int64) (*entity.User, error)
	GetPRsForReviewer(userID int64) ([]entity.PullRequest, error)
}

type PRRepository interface {
	CreatePR(pr *entity.PullRequest) error
	GetPR(id int64) (*entity.PullRequest, error)
	UpdatePR(pr *entity.PullRequest) error
	GetActiveReviewers(teamName string, exclude []int64) ([]entity.User, error)
}
