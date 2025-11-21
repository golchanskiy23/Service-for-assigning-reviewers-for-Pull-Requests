package postgres

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) UpdateActive(userID string, active bool) (*entity.User, error) {
	_, err := r.db.Exec(
		`UPDATE users SET is_active = $1 WHERE user_id = $2`,
		active, userID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetUser(userID)
}

func (r *UserRepo) GetUser(userID string) (*entity.User, error) {
	u := entity.User{}
	err := r.db.QueryRow(`
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE user_id = $1`,
		userID,
	).Scan(&u.ID, &u.Username, &u.TeamName, &u.Active)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetPRsForReviewer(userID string) ([]entity.PullRequest, error) {
	rows, err := r.db.Query(`
        SELECT pr.pull_request_id, pr.pull_request_name,
               pr.author_id, pr.status, pr.created_at, pr.merged_at
        FROM pr_reviewers r
        JOIN pull_requests pr ON pr.pull_request_id = r.pull_request_id
        WHERE r.reviewer_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := []entity.PullRequest{}
	for rows.Next() {
		var pr entity.PullRequest
		rows.Scan(
			&pr.ID, &pr.Name, &pr.AuthorID,
			&pr.Status, &pr.CreatedAt, &pr.MergedAt,
		)
		prs = append(prs, pr)
	}
	return prs, nil
}
