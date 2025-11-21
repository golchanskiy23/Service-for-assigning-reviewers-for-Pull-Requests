package postgres

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
	"context"
)

type User struct {
	ID       int
	Name     string
	IsActive bool
}

type UserRepository interface {
	SetIsActive(ctx context.Context, userID int, active bool) error
	GetReview(ctx context.Context, userID int) ([]string, error)
}

type userPGRepository struct {
	db *database.DatabaseSource
}

func NewUserPGRepository(db *database.DatabaseSource) UserRepository {
	return &userPGRepository{db: db}
}

func (r *userPGRepository) SetIsActive(ctx context.Context, userID int, active bool) error {
	_, err := r.db.Pool.Exec(ctx,
		`UPDATE users SET is_active=$1 WHERE id=$2`, active, userID)
	return err
}

func (r *userPGRepository) GetReview(ctx context.Context, userID int) ([]string, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT review_text FROM reviews WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []string
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return nil, err
		}
		reviews = append(reviews, text)
	}
	return reviews, nil
}

/*func (r *userPGRepository) UpdateActive(userID string, active bool) (*entity.User, error) {
	_, err := r.db.Exec(
		`UPDATE users SET is_active = $1 WHERE user_id = $2`,
		active, userID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetUser(userID)
	return nil,nil
}

func (r *userPGRepository) GetUser(userID string) (*entity.User, error) {
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
*/
