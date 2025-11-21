package postgres

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
)

type Repository struct {
	Teams        TeamRepository
	Users        UserRepository
	PullRequests PullRequestRepository
}

func CreateNewDBRepository(db *database.DatabaseSource) *Repository {
	return &Repository{
		Teams:        NewTeamPGRepository(db),
		Users:        NewUserPGRepository(db),
		PullRequests: NewPullRequestPGRepository(db),
	}
}
