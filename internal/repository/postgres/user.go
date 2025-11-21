package postgres

import (
	"context"
	"errors"

	"Service-for-assigning-reviewers-for-Pull-Requests/internal/entity"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
)

type UserRepository interface {
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	SetIsActive(ctx context.Context, userID string, active bool) error
	GetActiveUsersByTeam(
		ctx context.Context,
		teamName string,
		exclude []string,
	) ([]*entity.User, error)
	//nolint:revive // func
	GetPRsForReviewer(ctx context.Context, userID string) ([]*entity.PullRequestShort, error)
	// MassDeactivateAndReassign deactivates given users (must belong to the same team)
	// and for every OPEN PR where they are reviewers removes them and tries
	// to reassign at least one active reviewer from the same team.
	MassDeactivateAndReassign(ctx context.Context, teamName string, userIDs []string) error
}

type userPGRepository struct {
	db *database.DatabaseSource
}

func NewUserPGRepository(db *database.DatabaseSource) UserRepository {
	return &userPGRepository{db: db}
}

func (r *userPGRepository) GetUser(
	ctx context.Context,
	userID string,
) (*entity.User, error) {
	var user entity.User

	err := r.db.Pool.QueryRow(ctx,
		`SELECT user_id, username, team_name, is_active
		 FROM users
		 WHERE user_id = $1`,
		userID,
	).Scan(
		&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)
	if err != nil {
		return nil, errors.New("NOT_FOUND")
	}

	return &user, nil
}

func (r *userPGRepository) SetIsActive(
	ctx context.Context,
	userID string,
	active bool,
) error {
	result, err := r.db.Pool.Exec(ctx,
		`UPDATE users SET is_active = $1 WHERE user_id = $2`, active, userID)
	if err != nil {
		return err
	}

	const noRowsAffected = 0
	if result.RowsAffected() == noRowsAffected {
		return errors.New("NOT_FOUND")
	}

	return nil
}

func (r *userPGRepository) GetActiveUsersByTeam(
	ctx context.Context,
	teamName string,
	exclude []string,
) ([]*entity.User, error) {
	query := `SELECT user_id, username, team_name, is_active
			  FROM users
			  WHERE team_name = $1 AND is_active = TRUE`
	args := []interface{}{teamName}

	const emptySlice = 0
	if len(exclude) > emptySlice {
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
		if err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.TeamName,
			&user.IsActive,
		); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *userPGRepository) GetPRsForReviewer(
	ctx context.Context,
	userID string,
) ([]*entity.PullRequestShort, error) {
	rows, err := r.db.Pool.Query(ctx,
		//nolint:revive // sql query
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
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
		); err != nil {
			return nil, err
		}

		prs = append(prs, &pr)
	}

	return prs, nil
}

// MassDeactivateAndReassign implements bulk deactivation and safe reassignment.
func (r *userPGRepository) MassDeactivateAndReassign(
	ctx context.Context,
	teamName string,
	userIDs []string,
) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	//nolint:errcheck // best-effort
	defer tx.Rollback(ctx)

	// Deactivate provided users (scoped to team for safety)
	if _, err := tx.Exec(ctx,
		`UPDATE users SET is_active = FALSE WHERE user_id = ANY($1::text[]) AND team_name = $2`,
		userIDs, teamName); err != nil {
		return err
	}

	// Find affected OPEN PRs that had any of these users as reviewers and get current reviewer list
	rows, err := tx.Query(ctx,
		`SELECT pr.pull_request_id, array_agg(prr.reviewer_id ORDER BY prr.assigned_at)
		 FROM pr_reviewers prr
		 JOIN pull_requests pr ON pr.pull_request_id = prr.pull_request_id
		 WHERE pr.status = 'OPEN' AND prr.reviewer_id = ANY($1::text[])
		 GROUP BY pr.pull_request_id`,
		userIDs,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var prID string
		var reviewers []string
		if err := rows.Scan(&prID, &reviewers); err != nil {
			return err
		}

		// Build remaining reviewers by excluding deactivated ones
		remaining := make([]string, 0, len(reviewers))
		for _, rID := range reviewers {
			skip := false
			for _, d := range userIDs {
				if rID == d {
					skip = true
					break
				}
			}
			if !skip {
				remaining = append(remaining, rID)
			}
		}

		// Delete entries for deactivated users for this PR
		if _, err := tx.Exec(ctx,
			`DELETE FROM pr_reviewers WHERE pull_request_id = $1 AND reviewer_id = ANY($2::text[])`,
			prID, userIDs); err != nil {
			return err
		}

		// If no remaining reviewers, try to add one active team member
		if len(remaining) == 0 {
			// find one active candidate in the same team (exclude deactivated list)
			var candidate string
			err := tx.QueryRow(ctx,
				`SELECT user_id FROM users
				 WHERE team_name = $1 AND is_active = TRUE AND NOT (user_id = ANY($2::text[]))
				 LIMIT 1`,
				teamName, userIDs).Scan(&candidate)

			if err == nil && candidate != "" {
				if _, err := tx.Exec(ctx,
					`INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`,
					prID, candidate); err != nil {
					return err
				}
			}
			// if no candidate found, leave PR without reviewers — caller may handle this
		}
	}

	return tx.Commit(ctx)
}
