package postgres

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
	"context"
	"errors"
)

type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	SetIsActive(ctx context.Context, userID string, active bool) error
	GetActiveUsersByTeam(ctx context.Context, teamName string, exclude []string) ([]*entity.User, error)
	GetPRsForReviewer(ctx context.Context, userID string) ([]*entity.PullRequestShort, error)
}

type userPGRepository struct {
	db *database.DatabaseSource
}

func NewUserPGRepository(db *database.DatabaseSource) UserRepository {
	return &userPGRepository{db: db}
}

func (r *userPGRepository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	var user entity.User
	err := r.db.Pool.QueryRow(ctx,
		`SELECT user_id, username, team_name, is_active
		 FROM users
		 WHERE user_id = $1`,
		userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		return nil, errors.New("NOT_FOUND")
	}
	return &user, nil
}

func (r *userPGRepository) SetIsActive(ctx context.Context, userID string, active bool) error {
	result, err := r.db.Pool.Exec(ctx,
		`UPDATE users SET is_active = $1 WHERE user_id = $2`, active, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("NOT_FOUND")
	}
	return nil
}

func (r *userPGRepository) GetActiveUsersByTeam(ctx context.Context, teamName string, exclude []string) ([]*entity.User, error) {
	query := `SELECT user_id, username, team_name, is_active
			  FROM users
			  WHERE team_name = $1 AND is_active = TRUE`
	args := []interface{}{teamName}

	if len(exclude) > 0 {
		query += ` AND NOT (user_id = ANY($2::text[]))`
		args = append(args, exclude)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userPGRepository) GetPRsForReviewer(ctx context.Context, userID string) ([]*entity.PullRequestShort, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		 FROM pr_reviewers prr
		 JOIN pull_requests pr ON pr.pull_request_id = prr.pull_request_id
		 WHERE prr.reviewer_id = $1
		 ORDER BY pr.created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*entity.PullRequestShort
	for rows.Next() {
		var pr entity.PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, nil
}
